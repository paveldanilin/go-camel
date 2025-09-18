package camel

import (
	"fmt"
)

// setErrorProcessor sets a camel.Exchange error
type setErrorProcessor struct {
	// stepName is a logical name of current operation.
	stepName string
	err      error
}

func newSetErrorProcessor(err error) *setErrorProcessor {
	return &setErrorProcessor{
		stepName: fmt.Sprintf("setError{err=%v}", err),
		err:      err,
	}
}

func (p *setErrorProcessor) WithStepName(stepName string) *setErrorProcessor {
	p.stepName = stepName
	return p
}

func (p *setErrorProcessor) Process(exchange *Exchange) {
	if !exchange.On(p.stepName) {
		return
	}

	exchange.SetError(p.err)
}
