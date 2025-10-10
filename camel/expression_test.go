package camel

import (
	"testing"
)

func TestSimpleExpression_Eq(t *testing.T) {
	e, err := newSimpleExpression("headers.a == 1")
	if err != nil {
		panic(err)
	}

	m := NewExchange(nil, nil)
	m.Message().SetHeader("a", 1)

	ret, err := e.eval(m)
	if err != nil {
		panic(err)
	}

	if ret.(bool) != true {
		t.Error("TestSimpleExpressionEq() = FALSE; want TRUE")
	}
}
