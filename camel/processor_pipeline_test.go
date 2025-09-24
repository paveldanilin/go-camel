package camel

import "testing"

func TestPipelineProcessor(t *testing.T) {
	sum := newPipelineProcessor("sum", false).
		addProcessor(newSetHeaderProcessor("set a", "a", newConstExpr(1))).
		addProcessor(newSetHeaderProcessor("set b", "b", newConstExpr(1))).
		addProcessor(newFuncProcessor("calc", func(exchange *Exchange) {

			a, _ := exchange.Message().Header("a")
			b, _ := exchange.Message().Header("b")
			exchange.Message().Body = a.(int) + b.(int)

		}))

	exchange := NewExchange(nil, nil)

	sum.Process(exchange)
	if exchange.IsError() {
		panic(exchange.Error)
	}

	result := exchange.Message().Body
	expected := 2

	if result != expected {
		t.Errorf("TestPipelineProcessor() = %d; want %d", result, expected)
	}
}
