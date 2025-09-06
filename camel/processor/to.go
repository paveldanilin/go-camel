package processor

import (
	"fmt"
	"github.com/paveldanilin/go-camel/camel"
)

type ToProcessor struct {
	// stepName is a logical name of current operation.
	stepName string
	uri      string
}

func To(uri string) *ToProcessor {
	return &ToProcessor{
		stepName: fmt.Sprintf("to{uri=%s}", uri),
		uri:      uri,
	}
}

func (p *ToProcessor) WithStepName(stepName string) *ToProcessor {
	p.stepName = stepName
	return p
}

func (p *ToProcessor) Process(exchange *camel.Exchange) {
	if !exchange.On(p.stepName) {
		return
	}

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
