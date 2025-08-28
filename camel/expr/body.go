package expr

import "github.com/paveldanilin/go-camel/camel"

type BodyExpr struct {
}

func Body() *BodyExpr {
	return &BodyExpr{}
}

func (e *BodyExpr) Eval(message *camel.Message) (any, error) {

	return message.Body, nil
}
