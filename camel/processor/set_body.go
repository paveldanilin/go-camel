package processor

import (
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
		value: value,
	}
}

func (p *SetBodyProcessor) SetStepName(stepName string) *SetBodyProcessor {
	p.stepName = stepName
	return p
}

func (p *SetBodyProcessor) Process(exchange *camel.Exchange) {
	exchange.PushStep(p.stepName)

	if err := exchange.CheckCancelOrTimeout(); err != nil {
		exchange.Error = err
		return
	}

	value, err := p.value.Eval(exchange)
	if err != nil {
		exchange.Error = err
		return
	}

	exchange.Message().Body = value
}
