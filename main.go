package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"sync/atomic"

	"github.com/google/uuid"

	ce "github.com/cloudevents/sdk-go"
	"github.com/cloudevents/sdk-go/pkg/binding"
	"knative.dev/eventing/pkg/kncloudevents"

	"github.com/cloudevents/sdk-go/pkg/bindings/kafka_sarama"
)

var (
	clientID         string
	config           *Config
	sender           *kafka_sarama.Sender
	ensureTopicMutex = &sync.Mutex{}

	// counters
	messagesCounter *int64 = new(int64)
)

func main() {

	var err error

	// configuration
	config, err = GetEnvConfig()
	if err != nil {
		log.Fatalf("Wrong configuration: %v", err)
		return
	}

	// setup kafka producer
	sender, err = createSender(config)
	if err != nil {
		return
	}

	// cloudevents client
	fmt.Println("Creating cloudevents default client")
	c, err := kncloudevents.NewDefaultClient()
	if err != nil {
		log.Fatal("Failed to create client, ", err)
		return
	}
	clientID = uuid.New().String()
	fmt.Printf("[%s] Cloudevents default client created\n", clientID)

	// start receiver
	errCh := make(chan error, 1)
	go func(c ce.Client, handler handler) {
		if err := c.StartReceiver(context.Background(), handler); err != nil {
			errCh <- err
		}
	}(c, receive)

	// wait until error
	select {
	case err := <-errCh:
		log.Fatalf("Failed to run HTTP server: %s\n", err)
	}
}

type handler func(ev ce.Event)

func receive(ev ce.Event) {
	inferenceServiceName := ev.Extensions()[InferenceServiceName].(string)
	fmt.Printf("[%s] CloudEvent received from %s\n", clientID, inferenceServiceName)

	// enrich event
	msg, err := NewMessage(ev)
	if err != nil {
		return
	}

	// get event binding
	em := binding.EventMessage(*msg)

	// send event to kafka
	go func(em binding.EventMessage) {
		if err := sender.Send(context.Background(), em); err != nil {
			log.Fatalf("[%s] Cannot produce kafka message [%s]: %v", clientID, ev.String(), err.Error())
		} else {
			messagesCount := atomic.AddInt64(messagesCounter, 1)
			fmt.Printf("[%s] Message from %s logged into %s. Total messages: %d\n", clientID, inferenceServiceName, config.KafkaTopic, messagesCount)
		}
	}(em)
}

func createSender(config *Config) (*kafka_sarama.Sender, error) {
	// kafka connection
	conn, err := GetKafkaConnection(config)
	if err != nil {
		log.Fatalf("Cannot create Kafka client to servers [%s]: %v", config.KafkaBrokers, err.Error())
		return nil, err
	}

	// ensure topic exists
	ensureTopicMutex.Lock() // avoid error: Cannot create topic: kafka server: Topic with this name already exists
	err = conn.EnsureTopic(config)
	ensureTopicMutex.Unlock()
	if err != nil {
		log.Fatalf("Cannot create topic [%s]: %v", config.KafkaTopic, err.Error())
		return nil, err
	}

	// create kafka producer
	prod, err := conn.GetProducer(config.KafkaTopic)
	if err != nil {
		log.Fatalf("Cannot create Kafka producer with topic [%s]: %v", config.KafkaTopic, err.Error())
		return nil, err
	}
	return &prod, nil
}
