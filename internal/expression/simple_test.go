package expression

import (
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
	"testing"
)

func TestSimple_Eval(t *testing.T) {
	e, err := NewSimple("header.a == 1")
	if err != nil {
		panic(err)
	}

	m := exchange.NewExchange(nil)
	m.Message().SetHeader("a", 1)

	ret, err := e.Eval(m)
	if err != nil {
		panic(err)
	}

	if ret.(bool) != true {
		t.Error("TestSimple_Eval() = FALSE; want TRUE")
	}
}
