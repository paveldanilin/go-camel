package camel

type messageBodyExpr struct {
}

func newMessageBodyExpr() *messageBodyExpr {
	return &messageBodyExpr{}
}

func (e *messageBodyExpr) Eval(exchange *Exchange) (any, error) {
	return exchange.Message().Body, nil
}
