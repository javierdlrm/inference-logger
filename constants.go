package main

// Config

const KafkaTopicEnvVar string = "KAFKA_TOPIC"
const KakfaBrokersEnvVar string = "KAFKA_BROKERS"

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

const Request = "request"
const Response = "response"

// Kafka Topic

const DefaultPartitions = 1
const DefaultReplicationFactor = 1

// DefaultTopicName returns a default name for Kafka topics
func DefaultTopicName(name string) string {
	return name + "-inference-topic"
}
