package main

// Config

const KafkaTopicEnvVar string = "KAFKA_TOPIC"
const KakfaBrokersEnvVar string = "KAFKA_BROKERS"
const KafkaTopicPartitionsEnvVar string = "KAFKA_TOPIC_PARTITIONS"
const KafkaTopicReplicationFactorEnvVar string = "KAFKA_TOPIC_REPLICATION_FACTOR"

// Defaults

const DefaultKafkaTopicPartitions int32 = 1
const DefaultKafkaTopicReplicationFactor int16 = 1

// CloudEvent attributes

const KFServingRequestType = "org.kubeflow.serving.inference.request"
const KFServingResponseType = "org.kubeflow.serving.inference.response"

const Type = "type"
const ID = "id"
const Time = "time"
const Payload = "payload"
const InferenceServiceName = "inferenceservicename"
const Endpoint = "endpoint"
const Namespace = "namespace"

const RequestType = "request"
const ResponseType = "response"

// DefaultTopicName returns a default name for Kafka topics
func DefaultTopicName(name string) string {
	return name + "-inference-topic"
}
