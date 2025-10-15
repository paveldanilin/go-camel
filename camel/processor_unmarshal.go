package camel

import (
	"fmt"
)

type unmarshalProcessor struct {
	name       string
	model      any
	dataFormat DataFormat
}

func newUnmarshalProcessor(name string, model any, dataFormat DataFormat) *unmarshalProcessor {
	return &unmarshalProcessor{
		name:       name,
		model:      model,
		dataFormat: dataFormat,
	}
}

func (p *unmarshalProcessor) getName() string {
	return p.name
}

func (p *unmarshalProcessor) Process(exchange *Exchange) {
	// Check Message.Body datatype
	var data []byte
	switch t := exchange.Message().Body.(type) {
	case string:
		data = []byte(t)
	case []byte:
		data = t
	default:
		exchange.SetError(fmt.Errorf("unmarshal: expected json data in message body, but got %T", exchange.Message().Body))
		return
	}

	body, err := p.dataFormat.Unmarshal(data, p.model)
	if err != nil {
		exchange.SetError(err)
		return
	}
	exchange.Message().Body = body
}
