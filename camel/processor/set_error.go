package processor

import (
	"github.com/paveldanilin/go-camel/camel"
)

// SetErrorProcessor sets a camel.Exchange error
type SetErrorProcessor struct {
	err error
}

func SetError(err error) *SetErrorProcessor {
	return &SetErrorProcessor{
		err: err,
	}
}

func (p *SetErrorProcessor) Process(exchange *camel.Exchange) {
	if err := exchange.CheckCancelOrTimeout(); err != nil {
		exchange.Error = err
		return
	}

	exchange.Error = p.err
}
