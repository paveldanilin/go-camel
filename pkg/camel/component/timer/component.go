package timer

import (
	"github.com/paveldanilin/go-camel/pkg/camel/api"
	camelUri "github.com/paveldanilin/go-camel/pkg/camel/uri"
)

type Component struct {
	exchangeFactory api.ExchangeFactory
}

func NewComponent() *Component {
	return &Component{}
}

func (c *Component) Id() string {
	return "timer"
}

func (c *Component) CreateEndpoint(uri string) (api.Endpoint, error) {
	parsedUri, err := camelUri.Parse(uri, nil)
	if err != nil {
		return nil, err
	}

	return NewEndpoint(parsedUri, c)
}

func (c *Component) SetExchangeFactory(f api.ExchangeFactory) {
	c.exchangeFactory = f
}
