package processor

import (
	"fmt"
	"github.com/paveldanilin/go-camel/camel"
)

// SetHeaderProcessor sets a camel.Message header
type SetHeaderProcessor struct {
	// stepName is a logical name of current operation.
	stepName string
	name     string
	value    camel.Expr
}

func SetHeader(name string, value camel.Expr) *SetHeaderProcessor {
	return &SetHeaderProcessor{
		stepName: fmt.Sprintf("setHeader{name=%s;value=%v}", name, value),
		name:     name,
		value:    value,
	}
}

func (p *SetHeaderProcessor) WithStepName(stepName string) *SetHeaderProcessor {
	p.stepName = stepName
	return p
}

func (p *SetHeaderProcessor) Process(exchange *camel.Exchange) {
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
