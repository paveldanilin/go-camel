package processor

import (
	"github.com/paveldanilin/go-camel/internal/expression"
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
)

type SetVariable struct {
	propertyName    string
	valueExpression expression.Expression
}

func NewSetVariable(name string, value expression.Expression) *SetVariable {
	return &SetVariable{
		propertyName:    name,
		valueExpression: value,
	}
}

func (p *SetVariable) Process(e *exchange.Exchange) {
	value, err := p.valueExpression.Eval(e)
	if err != nil {
		e.SetError(err)
		return
	}

	e.SetVariable(p.propertyName, value)
}
