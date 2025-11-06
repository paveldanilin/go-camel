package unmarshal

import (
	"fmt"
	"github.com/paveldanilin/go-camel/pkg/camel/api"
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
)

type unmarshalProcessor struct {
	routeName  string
	name       string
	model      any
	dataFormat api.DataFormat
}

func NewProcessor(routeName, name string, model any, dataFormat api.DataFormat) *unmarshalProcessor {
	return &unmarshalProcessor{
		routeName:  routeName,
		name:       name,
		model:      model,
		dataFormat: dataFormat,
	}
}

func (p *unmarshalProcessor) Name() string {
	return p.name
}

func (p *unmarshalProcessor) RouteName() string {
	return p.routeName
}

func (p *unmarshalProcessor) Process(e *exchange.Exchange) {
	// Check Message.Body datatype
	var data []byte
	switch t := e.Message().Body.(type) {
	case string:
		data = []byte(t)
	case []byte:
		data = t
	default:
		e.SetError(fmt.Errorf("unmarshal: expected json data in message body, but got %T", e.Message().Body))
		return
	}

	body, err := p.dataFormat.Unmarshal(data, p.model)
	if err != nil {
		e.SetError(err)
		return
	}
	e.Message().Body = body
}
