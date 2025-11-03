package direct

import (
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
)

type Producer struct {
	endpoint *Endpoint
}

func (p *Producer) Process(e *exchange.Exchange) {
	for _, producer := range p.endpoint.consumer.producers {
		producer.Process(e)
	}
}
