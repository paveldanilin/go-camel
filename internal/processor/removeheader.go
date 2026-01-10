package processor

import (
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
)

type RemoveHeader struct {
	headerNames []string
}

func NewRemoveHeader(names ...string) *RemoveHeader {
	return &RemoveHeader{
		headerNames: names,
	}
}

func (p *RemoveHeader) Process(e *exchange.Exchange) {
	for _, headerName := range p.headerNames {
		e.Message().RemoveHeader(headerName)
	}
}
