package processor

import (
	"fmt"
	"github.com/paveldanilin/go-camel/camel"
)

type ToProcessor struct {
	uri string
}

func To(uri string) *ToProcessor {
	return &ToProcessor{
		uri: uri,
	}
}

func (p *ToProcessor) Process(exchange *camel.Exchange) {
	if err := exchange.CheckCancelOrTimeout(); err != nil {
		exchange.Error = err
		return
	}

	endpoint := exchange.Runtime().Endpoint(p.uri)
	if endpoint == nil {
		exchange.Error = fmt.Errorf("endpoint not found '%s'", p.uri)
		return
	}

	producer, err := endpoint.CreateProducer()
	if err != nil {
		exchange.Error = err
		return
	}

	producer.Process(exchange)
}
