package processor

import (
	"github.com/paveldanilin/go-camel/camel"
)

// SetBodyProcessor sets a camel.Message body
type SetBodyProcessor struct {
	value camel.Expr
}

func SetBody(value camel.Expr) *SetBodyProcessor {
	return &SetBodyProcessor{
		value: value,
	}
}

func (p *SetBodyProcessor) Process(exchange *camel.Exchange) {
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
