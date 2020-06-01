package main

import (
	"github.com/Shopify/sarama"
)

// EnsureTopic ensures a given topic is created
func (k *KafkaConnection) EnsureTopic(topic string) error {

	// check if already exists
	topics, err := k.Client.Topics()
	if err != nil {
		return err
	} else if contains(topics, topic) {
		return nil
	}

	// ensure cluster admin client initiated
	k.EnsureClusterAdmin()

	detail := &sarama.TopicDetail{
		NumPartitions:     DefaultPartitions,
		ReplicationFactor: DefaultReplicationFactor,
	}
	return k.ClusterAdmin.CreateTopic(topic, detail, false)
}
