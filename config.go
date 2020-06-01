package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config defines logger configuration
type Config struct {
	KafkaBrokers                []string
	KafkaTopic                  string
	KafkaTopicPartitions        int32
	KafkaTopicReplicationFactor int16
}

// GetEnvConfig obtains Config from environment
func GetEnvConfig() (*Config, error) {
	topic := os.Getenv(KafkaTopicEnvVar)
	brokers, err := getKafkaBrokers()
	if err != nil {
		return nil, err
	}
	partitions, err := getKafkaTopicPartitions()
	if err != nil {
		return nil, err
	}
	replication, err := getKafkaTopicReplicationFactor()
	if err != nil {
		return nil, err
	}

	return &Config{
		KafkaBrokers:                brokers,
		KafkaTopic:                  topic,
		KafkaTopicPartitions:        int32(partitions),
		KafkaTopicReplicationFactor: int16(replication),
	}, nil
}

// EnsureKafkaTopicVar ensures KafkaTopic config is set
func (c *Config) EnsureKafkaTopicVar(inferenceService string) {
	if c.KafkaTopic == "" {
		c.KafkaTopic = DefaultTopicName(inferenceService)
	}
}

func getKafkaBrokers() ([]string, error) {
	brokers := os.Getenv(KakfaBrokersEnvVar)
	if brokers == "" {
		return nil, fmt.Errorf("Brokers env var is required")
	}
	return strings.Split(brokers, ","), nil
}

func getKafkaTopicPartitions() (int, error) {
	str := os.Getenv(KafkaTopicPartitionsEnvVar)
	if str == "" {
		return DefaultKafkaTopicPartitions, nil
	}
	return strconv.Atoi(str)
}

func getKafkaTopicReplicationFactor() (int, error) {
	str := os.Getenv(KafkaTopicReplicationFactorEnvVar)
	if str == "" {
		return DefaultKafkaTopicReplicationFactor, nil
	}
	return strconv.Atoi(str)
}
