package processor

import (
	"github.com/paveldanilin/go-camel/pkg/camel/api"
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
)

type To struct {
	producer api.Producer
}

func NewTo(producer api.Producer) *To {
	return &To{
		producer: producer,
	}
}

func (p *To) Process(e *exchange.Exchange) {
	p.producer.Process(e)
}
