package camel

import (
	"fmt"
)

type toProcessor struct {
	name string

	uri string
}

func newToProcessor(name, uri string) *toProcessor {
	return &toProcessor{
		name: name,
		uri:  uri,
	}
}

func (p *toProcessor) getName() string {
	return p.name
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
