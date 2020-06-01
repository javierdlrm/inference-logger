package main

import (
	"context"
	"fmt"
	"log"

	ce "github.com/cloudevents/sdk-go"
	"github.com/cloudevents/sdk-go/pkg/binding"
	"knative.dev/eventing/pkg/kncloudevents"
)

var config *Config

func main() {

	var err error

	// configuration
	config, err = GetEnvConfig()
	if err != nil {
		log.Fatalf("Wrong configuration: %v", err)
		return
	}

	// cloudevents client
	fmt.Println("Creating cloudevents default client")
	c, err := kncloudevents.NewDefaultClient()
	if err != nil {
		log.Fatal("Failed to create client, ", err)
	} else {
		fmt.Println("Cloudevents default client created")
	}

	// start receiver
	log.Fatal(c.StartReceiver(context.Background(), receive))
}

func receive(ev ce.Event) {
	inferenceServiceName := ev.Extensions()[InferenceServiceName].(string)
	fmt.Println("CloudEvent received from " + inferenceServiceName)

	// kafka connection
	conn, err := GetKafkaConnection(config.KafkaBrokers)
	if err != nil {
		log.Fatalf("Cannot create Kafka client to servers [%s]: %v", config.KafkaBrokers, err.Error())
		return
	}

	// enrich event
	eev := EnrichEvent(ev)

	// ensure topic exists
	config.EnsureKafkaTopicVar(inferenceServiceName)
	if err := conn.EnsureTopic(config.KafkaTopic); err != nil {
		log.Fatalf("Cannot create topic [%s]: %v", config.KafkaTopic, err.Error())
		return
	}

	// create kafka producer
	prod, err := conn.GetProducer(config.KafkaTopic)
	if err != nil {
		log.Fatalf("Cannot create Kafka producer with topic [%s]: %v", config.KafkaTopic, err.Error())
		return
	}

	// get event binding
	em := binding.EventMessage(eev)

	// send event to kafka
	if err := prod.Send(context.Background(), em); err != nil {
		log.Fatalf("Cannot produce kafka message [%s]: %v", ev.String(), err.Error())
	} else {
		fmt.Printf("Message from '%s' log into '%s'", inferenceServiceName, config.KafkaTopic)
	}
}
