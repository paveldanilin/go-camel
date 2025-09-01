package expr

import "github.com/paveldanilin/go-camel/camel"

type ConstExpr struct {
	value any
}

func Const(value any) *ConstExpr {
	return &ConstExpr{
		value: value,
	}
}

func (e *ConstExpr) Eval(_ *camel.Exchange) (any, error) {
	return e.value, nil
}
