package to

import (
	"github.com/paveldanilin/go-camel/pkg/camel/api"
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
)

type toProcessor struct {
	routeName string
	name      string
	uri       string
	endpoint  api.Endpoint
}

func NewProcessor(routeName, name, uri string, endpoint api.Endpoint) *toProcessor {
	return &toProcessor{
		routeName: routeName,
		name:      name,
		uri:       uri,
		endpoint:  endpoint,
	}
}

func (p *toProcessor) Name() string {
	return p.name
}

func (p *toProcessor) RouteName() string {
	return p.routeName
}

func (p *toProcessor) Process(e *exchange.Exchange) {
	producer, err := p.endpoint.CreateProducer()
	if err != nil {
		e.SetError(err)
		return
	}

	producer.Process(e)
}
