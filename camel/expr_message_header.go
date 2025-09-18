package camel

type messageHeaderExpr struct {
	name string
}

func newMessageHeaderExpr(name string) *messageHeaderExpr {
	return &messageHeaderExpr{
		name: name,
	}
}

func (e *messageHeaderExpr) Eval(exchange *Exchange) (any, error) {
	headerValue, _ := exchange.Message().Header(e.name)
	return headerValue, nil
}
