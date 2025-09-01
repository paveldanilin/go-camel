package expr

import "github.com/paveldanilin/go-camel/camel"

type Func func(exchange *camel.Exchange) (any, error)

func (e Func) Eval(exchange *camel.Exchange) (any, error) {
	return e(exchange)
}
