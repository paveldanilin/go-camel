package setproperty

import (
	"github.com/paveldanilin/go-camel/internal/expression"
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
)

type setPropertyProcessor struct {
	routeName       string
	name            string
	propertyName    string
	valueExpression expression.Expression
}

func NewProcessor(routeName, name, propertyName string, valueExpression expression.Expression) *setPropertyProcessor {
	return &setPropertyProcessor{
		routeName:       routeName,
		name:            name,
		propertyName:    propertyName,
		valueExpression: valueExpression,
	}
}

func (p *setPropertyProcessor) Name() string {
	return p.name
}

func (p *setPropertyProcessor) RouteName() string {
	return p.routeName
}

func (p *setPropertyProcessor) Process(e *exchange.Exchange) {
	value, err := p.valueExpression.Eval(e)
	if err != nil {
		e.SetError(err)
		return
	}

	e.SetProperty(p.propertyName, value)
}
