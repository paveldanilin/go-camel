package camel

import "testing"

func TestPipelineProcessor(t *testing.T) {
	sum := newPipelineProcessor().
		WithStepName("Sum").
		WithProcessor(newSetHeaderProcessor("a", newConstExpr(1)).WithStepName("Set 'a' argument")).
		WithProcessor(newSetHeaderProcessor("b", newConstExpr(1)).WithStepName("Set 'b' argument")).
		WithProcessor(funcProcessor(func(exchange *Exchange) {

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
