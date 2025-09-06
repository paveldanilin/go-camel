package processor

import (
	"fmt"
	"github.com/paveldanilin/go-camel/camel"
)

// SetErrorProcessor sets a camel.Exchange error
type SetErrorProcessor struct {
	// stepName is a logical name of current operation.
	stepName string
	err      error
}

func SetError(err error) *SetErrorProcessor {
	return &SetErrorProcessor{
		stepName: fmt.Sprintf("setError{err=%v}", err),
		err:      err,
	}
}

func (p *SetErrorProcessor) WithStepName(stepName string) *SetErrorProcessor {
	p.stepName = stepName
	return p
}

func (p *SetErrorProcessor) Process(exchange *camel.Exchange) {
	if !exchange.On(p.stepName) {
		return
	}

	exchange.SetError(p.err)
}
