package processor

import (
	"fmt"
	"github.com/paveldanilin/go-camel/pkg/camel/api"
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
)

type Unmarshal struct {
	model      any
	dataFormat api.DataFormat
}

func NewUnmarshal(model any, dataFormat api.DataFormat) *Unmarshal {
	return &Unmarshal{
		model:      model,
		dataFormat: dataFormat,
	}
}

func (p *Unmarshal) Process(e *exchange.Exchange) {
	// Check Message.Body datatype
	var data []byte
	switch t := e.Message().Body().(type) {
	case string:
		data = []byte(t)
	case []byte:
		data = t
	default:
		e.SetError(fmt.Errorf("unmarshal: expected json data in message body, but got %T", e.Message().Body()))
		return
	}

	body, err := p.dataFormat.Unmarshal(data, p.model)
	if err != nil {
		e.SetError(err)
		return
	}
	e.Message().SetBody(body)
}
