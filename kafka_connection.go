package main

import (
	"sync"

	"github.com/Shopify/sarama"
)

// KafkaConnection defines a connection with Kafka
type KafkaConnection struct {
	Brokers      []string
	Config       *sarama.Config
	Client       sarama.Client
	ClusterAdmin sarama.ClusterAdmin

	clusterAdminOnce sync.Once
}

// GetKafkaConnection returns a new KafkaConnection
func GetKafkaConnection(brokers []string) (*KafkaConnection, error) {
	config := getKafkaConfig()

	client, err := getKafkaClient(config, brokers)
	if err != nil {
		return nil, err
	}

	return &KafkaConnection{
		Brokers: brokers,
		Config:  config,
		Client:  client,
	}, nil
}

// EnsureClusterAdmin ensures ClusterAdmin client is created
func (k *KafkaConnection) EnsureClusterAdmin() error {
	if k.ClusterAdmin != nil {
		return nil
	}

	var err error
	k.ClusterAdmin, err = getKafkaClusterAdmin(k.Config, k.Brokers)
	return err
}

func getKafkaConfig() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V2_0_0_0
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	return config
}

func getKafkaClient(config *sarama.Config, brokers []string) (sarama.Client, error) {
	client, err := sarama.NewClient(brokers, config)
	return client, err
}

func getKafkaClusterAdmin(config *sarama.Config, brokers []string) (sarama.ClusterAdmin, error) {
	admin, err := sarama.NewClusterAdmin(brokers, config)
	if err != nil {
		return nil, err
	}
	return admin, nil
}
