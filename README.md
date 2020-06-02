# Inference-logger
Kubernetes service that writes inference data from KFServing logger into Kafka topics.

> This service is used by the [model-monitoring operator](https://github.com/javierdlrm/model-monitoring-operator)

## How to start?

Deploy the inference logger service specifying the kafka brokers as environmental variable (see the next section). The service will listen for incoming request at 'http://<name>.default/' where <name> is the name you chose for the service.

Create or update your InferenceService with KFServing serving your model and add the endpoint mentioned before to predictor:logger:url in the specification. (You might want to specify also predictor:logger:mode all).

Make some predictions with your model and see the inference logs (request and response) in the corresponding kafka topic.

> Scripts for deployments and simple tests can be found in [model-monitoring-deploy](https://github.com/javierdlrm/model-monitoring-deploy)

## How it works?

Once the service is deploy, a cloudevents client is created using the [sdk-go](https://github.com/cloudevents/sdk-go) and starts to listen for cloudevents. The incoming events are enriched with additional information such as the timestamp of the requests and response and forwarded to kafka using [sarama](https://github.com/Shopify/sarama) and the [kafka_sarama binding](https://github.com/cloudevents/sdk-go/tree/master/protocol/kafka_sarama/v2).

## Environmental variables

The service requires some configuration to be able to reach the kafka cluster.


| Name                          | Type   | Description                                                              | Required | Default                                          |
|-------------------------------|--------|--------------------------------------------------------------------------|----------|--------------------------------------------------|
| KAFKA_BROKERS                 | string | Comma-separated list of brokers                                          | X        |                                                  |
| KAFKA_TOPIC                   | string | Name                                                                     |          | ${kfservingInferenceServiceName}-inference-topic |
| KAFKA_TOPIC_PARTITIONS        | string | Partitions for the topic (only applied if topic creation needed)         |          | 1                                                |
| KAFKA_TOPIC_REPLICATIONFACTOR | string | Replication factor for the topic (only applied if topic creation needed) |          | 1                                                |