package setheader

import (
	"github.com/paveldanilin/go-camel/internal/expression"
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
)

type setHeaderProcessor struct {
	routeName       string
	name            string
	headerName      string
	valueExpression expression.Expression
}

func NewProcessor(routeName, name, headerName string, valueExpression expression.Expression) *setHeaderProcessor {
	return &setHeaderProcessor{
		routeName:       routeName,
		name:            name,
		headerName:      headerName,
		valueExpression: valueExpression,
	}
}

func (p *setHeaderProcessor) Name() string {
	return p.name
}

func (p *setHeaderProcessor) RouteName() string {
	return p.routeName
}

func (p *setHeaderProcessor) Process(e *exchange.Exchange) {
	value, err := p.valueExpression.Eval(e)
	if err != nil {
		e.SetError(err)
		return
	}

	e.Message().SetHeader(p.headerName, value)
}
