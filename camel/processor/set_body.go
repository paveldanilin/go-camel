package processor

import (
	"fmt"
	"github.com/paveldanilin/go-camel/camel"
)

// SetBodyProcessor sets a camel.Message body
type SetBodyProcessor struct {
	// stepName is a logical name of current operation.
	stepName string
	value    camel.Expr
}

func SetBody(value camel.Expr) *SetBodyProcessor {
	return &SetBodyProcessor{
		stepName: fmt.Sprintf("setBody{expr=%s}", value),
		value:    value,
	}
}

func (p *SetBodyProcessor) WithStepName(stepName string) *SetBodyProcessor {
	p.stepName = stepName
	return p
}

func (p *SetBodyProcessor) Process(exchange *camel.Exchange) {
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
