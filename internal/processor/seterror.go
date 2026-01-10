package processor

import (
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
)

type SetError struct {
	err error
}

func NewSetError(err error) *SetError {
	return &SetError{
		err: err,
	}
}

func (p *SetError) Process(e *exchange.Exchange) {
	e.SetError(p.err)
}
