package camel

import (
	"testing"
)

func TestSimpleExpressionEq(t *testing.T) {
	exprEq, err := newSimpleExpr("header.a == 1")
	if err != nil {
		panic(err)
	}

	m := NewExchange(nil, nil)
	m.Message().SetHeader("a", 1)

	ret, err := exprEq.Eval(m)
	if err != nil {
		panic(err)
	}

	if ret.(bool) != true {
		t.Error("TestSimpleExpressionEq() = FALSE; want TRUE")
	}
}
