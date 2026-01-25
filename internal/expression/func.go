package expression

import "github.com/paveldanilin/go-camel/pkg/camel/exchange"

type Func struct {
	userFunc func(e *exchange.Exchange) (any, error)
}

func NewFunc(userFunc func(e *exchange.Exchange) (any, error)) *Func {
	return &Func{userFunc: userFunc}
}

func (f *Func) Eval(e *exchange.Exchange) (any, error) {
	return f.userFunc(e)
}
