package main

import (
	"github.com/Shopify/sarama"
)

// EnsureTopic ensures a given topic is created
func (k *KafkaConnection) EnsureTopic(c *Config) error {

	// check if already exists
	topics, err := k.Client.Topics()
	if err != nil {
		return err
	} else if contains(topics, c.KafkaTopic) {
		return nil
	}

	// ensure cluster admin client initiated
	k.EnsureClusterAdmin()
	defer func() { _ = k.ClusterAdmin.Close() }()

	detail := &sarama.TopicDetail{
		NumPartitions:     c.KafkaTopicPartitions,
		ReplicationFactor: c.KafkaTopicReplicationFactor,
	}
	return k.ClusterAdmin.CreateTopic(c.KafkaTopic, detail, false)
}
