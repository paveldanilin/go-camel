package processor

import (
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
)

type SetBody struct {
	valueExpression Expression
}

func NewSetBody(value Expression) *SetBody {
	return &SetBody{
		valueExpression: value,
	}
}

func (p *SetBody) Process(e *exchange.Exchange) {
	value, err := p.valueExpression.Eval(e)
	if err != nil {
		e.SetError(err)
		return
	}

	e.SetMessageBody(value)
}
