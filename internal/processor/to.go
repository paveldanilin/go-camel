package processor

import (
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
)

type To struct {
	producer Processor
}

func NewTo(producer Processor) *To {
	return &To{
		producer: producer,
	}
}

func (p *To) Process(e *exchange.Exchange) {
	p.producer.Process(e)
}
