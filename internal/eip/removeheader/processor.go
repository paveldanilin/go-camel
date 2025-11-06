package removeheader

import (
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
)

type removeHeaderProcessor struct {
	routeName   string
	name        string
	headerNames []string
}

func NewProcessor(routeName, name string, headerNames ...string) *removeHeaderProcessor {
	return &removeHeaderProcessor{
		routeName:   routeName,
		name:        name,
		headerNames: headerNames,
	}
}

func (p *removeHeaderProcessor) Name() string {
	return p.name
}

func (p *removeHeaderProcessor) RouteName() string {
	return p.routeName
}

func (p *removeHeaderProcessor) Process(e *exchange.Exchange) {
	for _, headerName := range p.headerNames {
		e.Message().RemoveHeader(headerName)
	}
}
