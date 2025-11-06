package setbody

import (
	"github.com/paveldanilin/go-camel/internal/expression"
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
)

type setBodyProcessor struct {
	routeName       string
	name            string
	valueExpression expression.Expression
}

func NewProcessor(routeName, name string, valueExpression expression.Expression) *setBodyProcessor {
	return &setBodyProcessor{
		routeName:       routeName,
		name:            name,
		valueExpression: valueExpression,
	}
}

func (p *setBodyProcessor) Name() string {
	return p.name
}

func (p *setBodyProcessor) RouteName() string {
	return p.routeName
}

func (p *setBodyProcessor) Process(e *exchange.Exchange) {
	value, err := p.valueExpression.Eval(e)
	if err != nil {
		e.SetError(err)
		return
	}

	e.Message().Body = value
}
