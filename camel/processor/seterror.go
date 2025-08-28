package processor

import (
	"github.com/paveldanilin/go-camel/camel"
)

type SetErrorProcessor struct {
	err error
}

func SetError(err error) *SetErrorProcessor {

	return &SetErrorProcessor{
		err: err,
	}
}

func (p *SetErrorProcessor) Process(message *camel.Message) {

	message.Error = p.err
}
