package direct

import (
	"github.com/paveldanilin/go-camel/camel"
)

type Component struct {
}

func NewComponent() *Component {
	return &Component{}
}

func (c Component) Id() string {

	return "direct"
}

func (c Component) CreateEndpoint(uri string) (camel.Endpoint, error) {

	return &Endpoint{
		uri: uri,
	}, nil
}
