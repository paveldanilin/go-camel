package camel

type funcExpr func(exchange *Exchange) (any, error)

func (fn funcExpr) Eval(exchange *Exchange) (any, error) {
	return fn(exchange)
}
