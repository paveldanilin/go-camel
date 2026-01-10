package processor

import (
	"github.com/paveldanilin/go-camel/internal/expression"
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
)

type SetBody struct {
	valueExpression expression.Expression
}

func NewSetBody(valueExpression expression.Expression) *SetBody {
	return &SetBody{
		valueExpression: valueExpression,
	}
}

func (p *SetBody) Process(e *exchange.Exchange) {
	value, err := p.valueExpression.Eval(e)
	if err != nil {
		e.SetError(err)
		return
	}

	e.Message().SetBody(value)
}
