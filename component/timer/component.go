package timer

import (
	"github.com/paveldanilin/go-camel/camel"
)

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
	parsedUri, err := camel.Parse(uri, nil)
	if err != nil {
		return nil, err
	}

	return NewEndpoint(parsedUri, c)
}

func (c *Component) SetRuntime(rt *camel.Runtime) {
	c.runtime = rt
}
