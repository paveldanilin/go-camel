package processor

import (
	"github.com/paveldanilin/go-camel/pkg/camel/api"
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
)

type Marshal struct {
	dataFormat api.DataFormat
}

func NewMarshal(dataFormat api.DataFormat) *Marshal {
	return &Marshal{
		dataFormat: dataFormat,
	}
}

func (p *Marshal) Process(e *exchange.Exchange) {
	body, err := p.dataFormat.Marshal(e.Message().Body())
	if err != nil {
		e.SetError(err)
		return
	}
	e.Message().SetBody(body)
}
