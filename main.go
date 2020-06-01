package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	ce "github.com/cloudevents/sdk-go"
	"github.com/cloudevents/sdk-go/pkg/binding"
	"knative.dev/eventing/pkg/kncloudevents"
)

var (
	kfkTopic   string = os.Getenv(KAFKA_TOPIC)
	kfkBrokers string = os.Getenv(KAFKA_BROKERS)
)

func checkConfig() error {
	var vars []string

	if kfkTopic == "" {
		vars = append(vars, KAFKA_TOPIC)
	}
	if kfkBrokers == "" {
		vars = append(vars, KAFKA_BROKERS)
	}

	if len(vars) > 0 {
		return errors.New(fmt.Sprintf("Env vars must be set [%s]", strings.Join(vars, ", ")))
	}
	return nil
}

func receive(ev ce.Event) {
	fmt.Println("[InferenceLogger] CloudEvent received from " + fmt.Sprintf("%v", ev.Extensions()[InferenceServiceName]))

	// enrich event
	eev := enrichEvent(ev)

	// create kafka client
	kfkClient := getKafkaClient(strings.Split(kfkBrokers, ","))

	// create kafka producer
	kfkProducer := getKafkaProducer(kfkClient, kfkTopic)

	defer kfkProducer.Close(context.Background())

	// get event binding
	em := binding.EventMessage(eev)

	// send event to kafka
	if err := kfkProducer.Send(context.Background(), em); err != nil {
		log.Fatalf("[InferenceLogger] Cannot produce kafka message [%s]: %v", ev.String(), err.Error())
		fmt.Printf("[InferenceLogger] Cannot produce kafka message [%s]: %v", ev.String(), err.Error())
	} else {
		fmt.Println("[InferenceLogger] Message sent")
	}
}

func main() {

	// check configuration
	if err := checkConfig(); err != nil {
		log.Fatalf("[InferenceLogger] Wrong configuration: %v", err)
		fmt.Printf("[InferenceLogger] Wrong configuration: %v", err)
	}

	// create cloudevents client
	fmt.Println("[InferenceLogger] Creating cloudevents default client")
	c, err := kncloudevents.NewDefaultClient()
	if err != nil {
		log.Fatal("Failed to create client, ", err)
	} else {
		fmt.Println("[InferenceLogger] Cloudevents default client created")
	}

	// start receiver
	log.Fatal(c.StartReceiver(context.Background(), receive))
}
