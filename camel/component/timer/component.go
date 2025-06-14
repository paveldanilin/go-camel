package timer

import "github.com/paveldanilin/go-camel/camel"

type Component struct {
	context *camel.Context
}

func NewComponent() *Component {
	return &Component{}
}

func (c *Component) Id() string {
	return "timer"
}

func (c *Component) Endpoint(uri string) (camel.Endpoint, error) {
	return &Endpoint{
		component: c,
		uri:       "timer:" + uri,
	}, nil
}

func (c *Component) SetContext(context *camel.Context) {
	c.context = context
}
