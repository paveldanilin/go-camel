package marshal

import (
	"github.com/paveldanilin/go-camel/pkg/camel/api"
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
)

type marshalProcessor struct {
	routeName  string
	name       string
	dataFormat api.DataFormat
}

func NewProcessor(routeName, name string, dataFormat api.DataFormat) *marshalProcessor {
	return &marshalProcessor{
		routeName:  routeName,
		name:       name,
		dataFormat: dataFormat,
	}
}

func (p *marshalProcessor) Name() string {
	return p.name
}

func (p *marshalProcessor) RouteName() string {
	return p.routeName
}

func (p *marshalProcessor) Process(e *exchange.Exchange) {
	body, err := p.dataFormat.Marshal(e.Message().Body)
	if err != nil {
		e.SetError(err)
		return
	}
	e.Message().Body = body
}
