package processor

import (
	"fmt"
	"github.com/paveldanilin/go-camel/camel"
)

// SetPropertyProcessor sets camel.Exchange property
type SetPropertyProcessor struct {
	// stepName is a logical name of current operation.
	stepName string
	name     string
	value    camel.Expr
}

func SetProperty(name string, value camel.Expr) *SetPropertyProcessor {
	return &SetPropertyProcessor{
		stepName: fmt.Sprintf("setProperty{name=%s;value=%v}", name, value),
		name:     name,
		value:    value,
	}
}

func (p *SetPropertyProcessor) WithStepName(stepName string) *SetPropertyProcessor {
	p.stepName = stepName
	return p
}

func (p *SetPropertyProcessor) Process(exchange *camel.Exchange) {
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
