package fn

import (
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
)

type fnProcessor struct {
	routeName string
	name      string
	userFunc  func(*exchange.Exchange)
}

func NewProcessor(routeName, name string, userFunc func(*exchange.Exchange)) *fnProcessor {
	return &fnProcessor{
		routeName: routeName,
		name:      name,
		userFunc:  userFunc,
	}
}

func (p *fnProcessor) Name() string {
	return p.name
}

func (p *fnProcessor) RouteName() string {
	return p.routeName
}

func (p *fnProcessor) Process(e *exchange.Exchange) {
	p.userFunc(e)
}
