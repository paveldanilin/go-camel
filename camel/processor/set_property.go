package processor

import "github.com/paveldanilin/go-camel/camel"

// SetPropertyProcessor sets camel.Exchange property
type SetPropertyProcessor struct {
	name  string
	value camel.Expr
}

func SetProperty(name string, value camel.Expr) *SetPropertyProcessor {
	return &SetPropertyProcessor{
		name:  name,
		value: value,
	}
}

func (p *SetPropertyProcessor) Process(exchange *camel.Exchange) {
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
