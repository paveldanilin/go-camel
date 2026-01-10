package processor

import (
	"github.com/paveldanilin/go-camel/internal/expression"
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
)

type SetHeader struct {
	headerName      string
	valueExpression expression.Expression
}

func NewSetHeader(name string, value expression.Expression) *SetHeader {
	return &SetHeader{
		headerName:      name,
		valueExpression: value,
	}
}

func (p *SetHeader) Process(e *exchange.Exchange) {
	value, err := p.valueExpression.Eval(e)
	if err != nil {
		e.SetError(err)
		return
	}

	e.Message().SetHeader(p.headerName, value)
}
