package camel

import (
	"fmt"
)

type toProcessor struct {
	// stepName is a logical name of current operation.
	stepName string
	uri      string
}

func newToProcessor(uri string) *toProcessor {
	return &toProcessor{
		stepName: fmt.Sprintf("to{uri=%s}", uri),
		uri:      uri,
	}
}

func (p *toProcessor) WithStepName(stepName string) *toProcessor {
	p.stepName = stepName
	return p
}

func (p *toProcessor) Process(exchange *Exchange) {
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
