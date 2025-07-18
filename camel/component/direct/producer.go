package direct

import (
	"github.com/paveldanilin/go-camel/camel"
)

type Producer struct {
	endpoint *Endpoint
}

func (p *Producer) Process(message *camel.Message) {

	for _, producer := range p.endpoint.consumer.producers {
		producer.Process(message)
	}
}
