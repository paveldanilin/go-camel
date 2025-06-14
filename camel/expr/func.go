package expr

import "github.com/paveldanilin/go-camel/camel"

type Func func(message *camel.Message) (any, error)

func (e Func) Eval(message *camel.Message) (any, error) {

	return e(message)
}
