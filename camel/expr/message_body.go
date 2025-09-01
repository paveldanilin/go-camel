package expr

import "github.com/paveldanilin/go-camel/camel"

type MessageBodyExpr struct {
}

func MessageBody() *MessageBodyExpr {
	return &MessageBodyExpr{}
}

func (e *MessageBodyExpr) Eval(exchange *camel.Exchange) (any, error) {
	return exchange.Message().Body, nil
}
