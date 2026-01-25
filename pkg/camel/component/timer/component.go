package timer

import (
	"github.com/paveldanilin/go-camel/pkg/camel"
	camelUri "github.com/paveldanilin/go-camel/pkg/camel/uri"
)

type Component struct {
	exchangeFactory camel.ExchangeFactory
}

func NewComponent() *Component {
	return &Component{}
}

func (c *Component) ComponentId() string {
	return "timer"
}

func (c *Component) CreateEndpoint(uri string) (camel.Endpoint, error) {
	parsedUri, err := camelUri.Parse(uri, nil)
	if err != nil {
		return nil, err
	}

	return NewEndpoint(parsedUri, c)
}

func (c *Component) SetExchangeFactory(f camel.ExchangeFactory) {
	c.exchangeFactory = f
}
