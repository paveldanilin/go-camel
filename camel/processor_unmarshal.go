package camel

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type unmarshalProcessor struct {
	format string
	model  any
	data   []byte
}

func newUnmarshalProcessor(format string, model any, data []byte) unmarshalProcessor {
	return unmarshalProcessor{
		format: format,
		model:  model,
		data:   data,
	}
}

func (j DataFormatJson) Unmarshal(jsonData []byte, model any) (any, error) {
	modelType := reflect.TypeOf(model)
	newValuePtr := reflect.New(modelType)
	target := newValuePtr.Interface()

	if err := json.Unmarshal(jsonData, target); err != nil {
		return nil, fmt.Errorf("desearilization error json: %w", err)
	}

	return newValuePtr.Elem().Interface(), nil
}

func (j unmarshalProcessor) Process(exchange *Exchange) {
	var d DataFormat
	switch j.format {
	case "json":
		d = DataFormatJson{}
	default:
		fmt.Println("unknown data format")
	}
	model, err := d.Unmarshal(j.data, j.model)
	exchange.Message().Body = model
	exchange.err = err
}
