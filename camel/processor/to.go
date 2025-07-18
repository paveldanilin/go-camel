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

func (p *ToProcessor) Process(message *camel.Message) {

	producer, err := message.Runtime().Endpoint(p.uri).CreateProducer()
	if err != nil {
		message.SetError(err)
		return
	}

	producer.Process(message)
}
