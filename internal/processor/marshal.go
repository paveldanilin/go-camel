package processor

import (
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
)

type Marshal struct {
	dataFormat DataFormat
}

func NewMarshal(dataFormat DataFormat) *Marshal {
	return &Marshal{
		dataFormat: dataFormat,
	}
}

func (p *Marshal) Process(e *exchange.Exchange) {
	body, err := p.dataFormat.Marshal(e.MessageBody())
	if err != nil {
		e.SetError(err)
		return
	}
	e.SetMessageBody(body)
}
