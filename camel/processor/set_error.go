package processor

import (
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
		err: err,
	}
}

func (p *SetErrorProcessor) SetStepName(stepName string) *SetErrorProcessor {
	p.stepName = stepName
	return p
}

func (p *SetErrorProcessor) Process(exchange *camel.Exchange) {
	exchange.PushStep(p.stepName)

	if err := exchange.CheckCancelOrTimeout(); err != nil {
		exchange.Error = err
		return
	}

	exchange.Error = p.err
}
