package timer

import "github.com/paveldanilin/go-camel/camel"

type Component struct {
	runtime *camel.Runtime
}

func NewComponent() *Component {
	return &Component{}
}

func (c *Component) Id() string {

	return "timer"
}

func (c *Component) CreateEndpoint(uri string) (camel.Endpoint, error) {

	return &Endpoint{
		component: c,
		uri:       uri,
	}, nil
}

func (c *Component) SetRuntime(context *camel.Runtime) {

	c.runtime = context
}
