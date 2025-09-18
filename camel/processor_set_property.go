package camel

import (
	"fmt"
)

// setPropertyProcessor sets camel.Exchange property
type setPropertyProcessor struct {
	// stepName is a logical name of current operation.
	stepName string
	name     string
	value    Expr
}

func newSetPropertyProcessor(name string, value Expr) *setPropertyProcessor {
	return &setPropertyProcessor{
		stepName: fmt.Sprintf("setProperty{name=%s;value=%v}", name, value),
		name:     name,
		value:    value,
	}
}

func (p *setPropertyProcessor) WithStepName(stepName string) *setPropertyProcessor {
	p.stepName = stepName
	return p
}

func (p *setPropertyProcessor) Process(exchange *Exchange) {
	if !exchange.On(p.stepName) {
		return
	}

	value, err := p.value.Eval(exchange)
	if err != nil {
		exchange.SetError(err)
		return
	}

	exchange.SetProperty(p.name, value)
}
