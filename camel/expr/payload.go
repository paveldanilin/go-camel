package expr

import "github.com/paveldanilin/go-camel/camel"

type PayloadExpr struct {
}

func Payload() *PayloadExpr {
	return &PayloadExpr{}
}

func (e *PayloadExpr) Eval(message *camel.Message) (any, error) {

	return message.Payload(), nil
}
