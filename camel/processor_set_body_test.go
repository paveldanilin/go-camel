package camel

import "testing"

func TestSetBodyProcessor(t *testing.T) {
	mul := newSetBodyProcessor("set body", funcExpression(func(exchange *Exchange) (any, error) {

		a := exchange.Message().MustHeader("a")
		b := exchange.Message().MustHeader("b")

		return a.(int) * b.(int), nil
	}))

	exchange := NewExchange(nil, nil)
	exchange.Message().SetHeader("a", 2)
	exchange.Message().SetHeader("b", 3)

	mul.Process(exchange)
	if exchange.IsError() {
		panic(exchange.Error)
	}

	result := exchange.Message().Body
	expected := 6

	if result != expected {
		t.Errorf("TestSetBodyProcessor() = %d; want %d", result, expected)
	}
}
