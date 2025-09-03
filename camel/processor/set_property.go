package processor

import "github.com/paveldanilin/go-camel/camel"

// SetPropertyProcessor sets camel.Exchange property
type SetPropertyProcessor struct {
	// stepName is a logical name of current operation.
	stepName string
	name     string
	value    camel.Expr
}

func SetProperty(name string, value camel.Expr) *SetPropertyProcessor {
	return &SetPropertyProcessor{
		name:  name,
		value: value,
	}
}

func (p *SetPropertyProcessor) SetStepName(stepName string) *SetPropertyProcessor {
	p.stepName = stepName
	return p
}

func (p *SetPropertyProcessor) Process(exchange *camel.Exchange) {
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

	exchange.SetProperty(p.name, value)
}
