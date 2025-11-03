package timer

import (
	"github.com/paveldanilin/go-camel/pkg/camel/api"
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
	u "github.com/paveldanilin/go-camel/pkg/camel/uri"
)

type Component struct {
	exchangeFactory exchange.Factory
}

func NewComponent() *Component {
	return &Component{}
}

func (c *Component) Id() string {
	return "timer"
}

func (c *Component) CreateEndpoint(uri string) (api.Endpoint, error) {
	parsedUri, err := u.ParseURI(uri, nil)
	if err != nil {
		return nil, err
	}

	return NewEndpoint(parsedUri, c)
}

func (c *Component) SetExchangeFactory(f exchange.Factory) {
	c.exchangeFactory = f
}
