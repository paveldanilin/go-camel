package direct

import (
	"github.com/paveldanilin/go-camel/camel"
	u "github.com/paveldanilin/go-camel/camel/uri"
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
	parsedUri, err := u.Parse(uri, nil)
	if err != nil {
		return nil, err
	}

	return &Endpoint{
		uri: parsedUri,
	}, nil
}
