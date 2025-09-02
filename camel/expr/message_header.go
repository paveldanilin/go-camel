package expr

import "github.com/paveldanilin/go-camel/camel"

type MessageHeaderExpr struct {
	name string
}

func MessageHeader(name string) *MessageHeaderExpr {
	return &MessageHeaderExpr{
		name: name,
	}
}

func (e *MessageHeaderExpr) Eval(exchange *camel.Exchange) (any, error) {
	headerValue, _ := exchange.Message().Header(e.name)
	return headerValue, nil
}
