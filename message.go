package main

import (
	"encoding/json"
	"fmt"
	"sync/atomic"

	ce "github.com/cloudevents/sdk-go"
)

var (
	emptyRequestsCounter  *int64 = new(int64)
	emptyResponsesCounter *int64 = new(int64)
	requestsCounter       *int64 = new(int64)
	responsesCounter      *int64 = new(int64)
)

// Message schema used as kafka message
type Message struct {
	Type    string                 `json:"type"`
	ID      string                 `json:"id"`
	Time    int64                  `json:"time"`
	Payload map[string]interface{} `json:"payload"`
}

// NewMessage builds a message ready to log into Kafka
func NewMessage(ev ce.Event) (*ce.Event, error) {
	var (
		msg *Message
		err error

		msgType       string
		responsesIncr int64
		requestsIncr  int64
	)

	// check message type
	switch ev.Context.GetType() {
	case KFServingResponseType:
		responsesIncr, requestsIncr, msgType = 1, 0, ResponseType
	case KFServingRequestType:
		responsesIncr, requestsIncr, msgType = 0, 1, RequestType
	default:
		if err = fmt.Errorf("Unknown event type: %s", ev.Context.GetType()); err != nil {
			return nil, err
		}
	}

	// build message
	msg, err = buildMessage(msgType, ev.Context.GetID(), ev.Context.GetTime().Unix(), ev.Data)
	if err != nil {
		fmt.Printf("Error building message: %v", err)
		return nil, err
	}
	data, err := json.Marshal(msg)
	if err != nil {
		fmt.Printf("Error json.Marshal inference event map: %v", err)
		return nil, fmt.Errorf("Error json.Marshal inference event map: %v", err)
	}
	ev.SetData(data) // replace data with new inference event

	// status
	responsesCount := atomic.AddInt64(responsesCounter, responsesIncr)
	requestsCount := atomic.AddInt64(requestsCounter, requestsIncr)
	fmt.Printf("[%s] Total processed: %d responses, %d requests\n", clientID, responsesCount, requestsCount)

	return &ev, nil
}

func buildMessage(msgType string, id string, t int64, d interface{}) (*Message, error) {
	p, err := getDataAsMap(d)
	if err != nil {
		return nil, err
	}
	if err = validatePayload(msgType, p); err != nil {
		return nil, err
	}
	return &Message{
		Type:    msgType,
		ID:      id,
		Time:    t,
		Payload: p,
	}, nil
}

func validatePayload(pType string, p map[string]interface{}) error {
	if pType == RequestType {
		if ins := p["instances"]; ins == nil {
			emptyRequestsCount := atomic.AddInt64(emptyRequestsCounter, 1)
			fmt.Printf("[%s] Wrong payload: Request without instances. Total cases: %d \n", clientID, emptyRequestsCount)
			return fmt.Errorf("request without instances")
		}
	} else {
		if preds := p["predictions"]; preds == nil {
			emptyResponsesCount := atomic.AddInt64(emptyResponsesCounter, 1)
			fmt.Printf("[%s] Wrong payload: Response without predictions. Total cases: %d \n", clientID, emptyResponsesCount)
			return fmt.Errorf("response without predictions")
		}
	}
	return nil
}

func getDataAsMap(d interface{}) (map[string]interface{}, error) {
	data, ok := d.([]byte)
	if !ok {
		var err error
		data, err = json.Marshal(d)
		if err != nil {
			return nil, err
		}
	}
	var dataMap map[string]interface{}
	err := json.Unmarshal(data, &dataMap)
	if err != nil {
		fmt.Printf("Error json.Unmarshal event data: %v", err)
		return nil, err
	}
	return dataMap, nil
}

func getKeys(p map[string]interface{}) []string {
	keys, i := make([]string, len(p)), 0
	for k := range p {
		keys[i] = k
		i++
	}
	return keys
}
