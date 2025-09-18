package camel

import (
	"fmt"
)

// setHeaderProcessor sets a camel.Message header
type setHeaderProcessor struct {
	// stepName is a logical name of current operation.
	stepName string
	name     string
	value    Expr
}

func newSetHeaderProcessor(name string, value Expr) *setHeaderProcessor {
	return &setHeaderProcessor{
		stepName: fmt.Sprintf("setHeader{name=%s;value=%v}", name, value),
		name:     name,
		value:    value,
	}
}

func (p *setHeaderProcessor) WithStepName(stepName string) *setHeaderProcessor {
	p.stepName = stepName
	return p
}

func (p *setHeaderProcessor) Process(exchange *Exchange) {
	if !exchange.On(p.stepName) {
		return
	}

	value, err := p.value.Eval(exchange)
	if err != nil {
		exchange.SetError(err)
		return
	}

	exchange.Message().SetHeader(p.name, value)
}
