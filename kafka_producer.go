package main

import (
	"context"

	"github.com/cloudevents/sdk-go/pkg/bindings/kafka_sarama"
)

// GetProducer return a new Kafka producer
func (k *KafkaConnection) GetProducer(topic string) (kafka_sarama.Sender, error) {
	prod, err := kafka_sarama.NewSender(k.Client, topic)
	if err != nil {
		return *prod, err
	}
	defer prod.Close(context.Background())
	return *prod, nil
}
