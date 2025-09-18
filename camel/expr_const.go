package camel

type constExpr struct {
	value any
}

func newConstExpr(value any) *constExpr {
	return &constExpr{
		value: value,
	}
}

func (e *constExpr) Eval(_ *Exchange) (any, error) {
	return e.value, nil
}
