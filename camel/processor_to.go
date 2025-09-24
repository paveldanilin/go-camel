package camel

import (
	"fmt"
)

type toProcessor struct {
	id string

	uri string
}

func newToProcessor(id, uri string) *toProcessor {
	return &toProcessor{
		id:  id,
		uri: uri,
	}
}

func (p *toProcessor) getId() string {
	return p.id
}

func (p *toProcessor) Process(exchange *Exchange) {
	endpoint := exchange.Runtime().Endpoint(p.uri)
	if endpoint == nil {
		exchange.SetError(fmt.Errorf("endpoint not found '%s'", p.uri))
		return
	}

	producer, err := endpoint.CreateProducer()
	if err != nil {
		exchange.SetError(err)
		return
	}

	producer.Process(exchange)
}
