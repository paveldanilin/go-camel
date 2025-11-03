package setbody

import (
	"github.com/paveldanilin/go-camel/internal/expression"
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
	"testing"
)

func TestSetBodyProcessor(t *testing.T) {
	mul := NewProcessor("", "set body", expression.Func(func(e *exchange.Exchange) (any, error) {

		a := e.Message().MustHeader("a")
		b := e.Message().MustHeader("b")

		return a.(int) * b.(int), nil
	}))

	e := exchange.NewExchange(nil)
	e.Message().SetHeader("a", 2)
	e.Message().SetHeader("b", 3)

	mul.Process(e)
	if e.IsError() {
		panic(e.Error())
	}

	result := e.Message().Body
	expected := 6

	if result != expected {
		t.Errorf("TestSetBodyProcessor() = %d; want %d", result, expected)
	}
}
