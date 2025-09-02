package expr

import "github.com/paveldanilin/go-camel/camel"

type Func func(exchange *camel.Exchange) (any, error)

func (fn Func) Eval(exchange *camel.Exchange) (any, error) {
	return fn(exchange)
}
