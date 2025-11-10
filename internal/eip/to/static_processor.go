package to

import (
	"github.com/paveldanilin/go-camel/pkg/camel/api"
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
)

type staticToProcessor struct {
	routeName string
	name      string
	producer  api.Producer
}

func NewStaticProcessor(routeName, name string, producer api.Producer) *staticToProcessor {
	return &staticToProcessor{
		routeName: routeName,
		name:      name,
		producer:  producer,
	}
}

func (p *staticToProcessor) Name() string {
	return p.name
}

func (p *staticToProcessor) RouteName() string {
	return p.routeName
}

func (p *staticToProcessor) Process(e *exchange.Exchange) {
	p.producer.Process(e)
}
