package camel

import (
	"encoding/json"
	"fmt"
)

type marshalProcessor struct {
	format string
	model  any
}

func newMarshalProcessor(format string, model any, data []byte) marshalProcessor {
	return marshalProcessor{
		format: format,
		model:  model,
	}
}

func (j DataFormatJson) Marshal(model any) (string, error) {
	data, err := json.Marshal(model)
	var json string

	if err != nil {
		return "", fmt.Errorf("searilization error json: %w", err)
	} else {
		json = string(data)
	}

	return json, nil
}

func (j marshalProcessor) Process(exchange *Exchange) {
	var d DataFormat
	switch j.format {
	case "json":
		d = DataFormatJson{}
	default:
		fmt.Println("unknown data format")
	}
	model, err := d.Marshal(j.model)
	exchange.Message().Body = model
	exchange.err = err
}
