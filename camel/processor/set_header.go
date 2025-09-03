package processor

import "github.com/paveldanilin/go-camel/camel"

// SetHeaderProcessor sets a camel.Message header
type SetHeaderProcessor struct {
	// stepName is a logical name of current operation.
	stepName string
	name     string
	value    camel.Expr
}

func SetHeader(name string, value camel.Expr) *SetHeaderProcessor {
	return &SetHeaderProcessor{
		name:  name,
		value: value,
	}
}

func (p *SetHeaderProcessor) SetStepName(stepName string) *SetHeaderProcessor {
	p.stepName = stepName
	return p
}

func (p *SetHeaderProcessor) Process(exchange *camel.Exchange) {
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

	exchange.Message().SetHeader(p.name, value)
}
