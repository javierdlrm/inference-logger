package main

import (
	"encoding/json"
	"fmt"

	ce "github.com/cloudevents/sdk-go"
)

func enrichEvent(ev ce.Event) ce.Event {

	// inference type
	ievtype := Request
	if ev.Context.GetType() == KFServingResponseType {
		ievtype = Response
	}

	// data
	data, ok := ev.Data.([]byte)
	if !ok {
		var err error
		data, err = json.Marshal(ev.Data)
		if err != nil {
			data = []byte(err.Error())
		}
	}
	var dataMap map[string]interface{}
	err := json.Unmarshal(data, &dataMap)
	if err != nil {
		fmt.Printf("[InferenceLogger] Error json.Unmarshal event data: %v", err)
	}

	// inference event
	ievMap := map[string]interface{}{
		Type:    ievtype,
		ID:      ev.Context.GetID(),
		Time:    ev.Context.GetTime(),
		Payload: dataMap,
	}
	iev, err := json.Marshal(ievMap)
	if err != nil {
		fmt.Printf("[InferenceLogger] Error json.Marshal inference event map: %v", err)
	}

	// replace data with new inference event
	ev.SetData(iev)

	return ev
}
