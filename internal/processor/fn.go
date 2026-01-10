package processor

import (
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
)

type Fn struct {
	userFunc func(*exchange.Exchange)
}

func NewFn(userFunc func(*exchange.Exchange)) *Fn {
	return &Fn{
		userFunc: userFunc,
	}
}

func (p *Fn) Process(e *exchange.Exchange) {
	p.userFunc(e)
}
