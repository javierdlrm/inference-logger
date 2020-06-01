package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/cloudevents/sdk-go/pkg/bindings/kafka_sarama"
)

func getKafkaClient(brokers []string) sarama.Client {
	config := sarama.NewConfig()
	config.Version = sarama.V2_0_0_0
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	client, err := sarama.NewClient(brokers, config)
	if err != nil {
		brokersStr := strings.Join(brokers, ",")
		log.Fatalf("[InferenceLogger] Cannot create Kafka client to servers [%s]: %v", brokersStr, err.Error())
		fmt.Printf("[InferenceLogger] Cannot create Kafka client to servers [%s]: %v", brokersStr, err.Error())
	} else {
		fmt.Println("[InferenceLogger] Kafka client created")
	}

	return client
}

func getKafkaProducer(kfkClient sarama.Client, kfkTopic string) kafka_sarama.Sender {
	kfkProducer, err := kafka_sarama.NewSender(kfkClient, kfkTopic)
	if err != nil {
		log.Fatalf("[InferenceLogger] Cannot create Kafka producer with topic [%s]: %v", kfkTopic, err.Error())
		fmt.Printf("[InferenceLogger] Cannot create Kafka producer with topic [%s]: %v", kfkTopic, err.Error())
	} else {
		fmt.Println("[InferenceLogger] Kafka producer created")
	}

	return *kfkProducer
}
