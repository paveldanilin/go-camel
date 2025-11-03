package direct

import (
	"github.com/paveldanilin/go-camel/pkg/camel/api"
	u "github.com/paveldanilin/go-camel/pkg/camel/uri"
)

type Component struct {
}

func NewComponent() *Component {
	return &Component{}
}

func (c Component) Id() string {
	return "direct"
}

func (c Component) CreateEndpoint(uri string) (api.Endpoint, error) {
	parsedUri, err := u.ParseURI(uri, nil)
	if err != nil {
		return nil, err
	}

	return &Endpoint{
		uri:  parsedUri,
		name: parsedUri.Path(),
	}, nil
}
