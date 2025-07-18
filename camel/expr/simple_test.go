package expr

import (
	"github.com/paveldanilin/go-camel/camel"
	"testing"
)

func TestSimpleExpressionEq(t *testing.T) {

	exprEq, err := Simple("header.a == 1")
	if err != nil {
		panic(err)
	}

	m := camel.NewMessage()
	m.SetHeader("a", 1)

	ret, err := exprEq.Eval(m)
	if err != nil {
		panic(err)
	}

	if ret.(bool) != true {
		t.Error("TestSimpleExpressionEq() = FALSE; want TRUE")
	}
}
