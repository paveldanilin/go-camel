package camel

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type jsonProcessor struct {
	json      []byte
	model     any
	operation JSON_OPERATION
}
type JSON_OPERATION int

const (
	Marshal JSON_OPERATION = iota
	Unmarshal
)

func newJsonProcessor(operation JSON_OPERATION, model any, json []byte) jsonProcessor {
	return jsonProcessor{
		json:      json,
		model:     model,
		operation: operation,
	}
}

func unmarshal(jsonData []byte, t reflect.Type) (any, error) {
	newValuePtr := reflect.New(t)
	target := newValuePtr.Interface()

	if err := json.Unmarshal(jsonData, target); err != nil {
		return nil, fmt.Errorf("desearilization error json: %w", err)
	}

	return newValuePtr.Elem().Interface(), nil
}

func marshal(model any) (string, error) {
	data, err := json.Marshal(model)
	var json string

	if err != nil {
		return "", fmt.Errorf("searilization error json: %w", err)
	} else {
		json = string(data)
	}

	return json, nil
}

func (j jsonProcessor) Process(exchange *Exchange) {
	switch j.operation {
	case Marshal:
		json, err := marshal(j.model)
		exchange.Message().Body = json
		exchange.err = err
	case Unmarshal:
		model, err := unmarshal([]byte(j.json), reflect.TypeOf(j.model))
		exchange.Message().Body = model
		exchange.err = err
	default:
		fmt.Println("operation should be marshal or unmarshal")
	}
}
