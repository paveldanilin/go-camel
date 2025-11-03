package seterror

import (
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
)

type setErrorProcessor struct {
	routeName string
	name      string
	err       error
}

func NewProcessor(routeName, name string, err error) *setErrorProcessor {
	return &setErrorProcessor{
		routeName: routeName,
		name:      name,
		err:       err,
	}
}

func (p *setErrorProcessor) Name() string {
	return p.name
}

func (p *setErrorProcessor) RouteName() string {
	return p.routeName
}

func (p *setErrorProcessor) Process(e *exchange.Exchange) {
	e.SetError(p.err)
}
