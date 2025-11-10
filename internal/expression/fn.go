package expression

import "github.com/paveldanilin/go-camel/pkg/camel/exchange"

type fn struct {
	userFunc func(e *exchange.Exchange) (any, error)
}

func NewFunc(userFunc func(e *exchange.Exchange) (any, error)) *fn {
	return &fn{userFunc: userFunc}
}

func (f *fn) Eval(e *exchange.Exchange) (any, error) {
	return f.userFunc(e)
}
