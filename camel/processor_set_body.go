package camel

import (
	"fmt"
)

// setBodyProcessor sets a camel.Message body
type setBodyProcessor struct {
	// stepName is a logical name of current operation.
	stepName string
	value    Expr
}

func newSetBodyProcessor(value Expr) *setBodyProcessor {
	return &setBodyProcessor{
		stepName: fmt.Sprintf("setBody{expr=%s}", value),
		value:    value,
	}
}

func (p *setBodyProcessor) WithStepName(stepName string) *setBodyProcessor {
	p.stepName = stepName
	return p
}

func (p *setBodyProcessor) Process(exchange *Exchange) {
	if !exchange.On(p.stepName) {
		return
	}

	value, err := p.value.Eval(exchange)
	if err != nil {
		exchange.SetError(err)
		return
	}

	exchange.Message().Body = value
}
