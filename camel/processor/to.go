package processor

import (
	"github.com/paveldanilin/go-camel/camel"
)

type ToProcessor struct {
	uri string // TODO: implement uri
}

func To(uri string) *ToProcessor {
	return &ToProcessor{
		uri: uri,
	}
}

func (p *ToProcessor) Process(message *camel.Message) error {

	pp, err := message.Context().Endpoint(p.uri).CreateProducer()
	if err != nil {
		return err
	}

	return pp.Process(message)
}
