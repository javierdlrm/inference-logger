package main

import (
	"fmt"
	"os"
	"strings"
)

// Config defines logger configuration
type Config struct {
	KafkaBrokers []string
	KafkaTopic   string
}

// GetEnvConfig obtains Config from environment
func GetEnvConfig() (*Config, error) {
	topic := os.Getenv(KafkaTopicEnvVar)
	brokers := os.Getenv(KakfaBrokersEnvVar)

	if brokers == "" {
		return nil, fmt.Errorf("Brokers env var must be set")
	}
	return &Config{
		KafkaBrokers: strings.Split(brokers, ","),
		KafkaTopic:   topic,
	}, nil
}

// EnsureKafkaTopicVar ensures KafkaTopic config is set
func (c *Config) EnsureKafkaTopicVar(inferenceService string) {
	if c.KafkaTopic == "" {
		c.KafkaTopic = DefaultTopicName(inferenceService)
	}
}
