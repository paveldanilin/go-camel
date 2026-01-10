package processor

import (
	"github.com/paveldanilin/go-camel/internal/expression"
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
)

type SetProperty struct {
	propertyName    string
	valueExpression expression.Expression
}

func NewSetProperty(name string, value expression.Expression) *SetProperty {
	return &SetProperty{
		propertyName:    name,
		valueExpression: value,
	}
}

func (p *SetProperty) Process(e *exchange.Exchange) {
	value, err := p.valueExpression.Eval(e)
	if err != nil {
		e.SetError(err)
		return
	}

	e.SetProperty(p.propertyName, value)
}
