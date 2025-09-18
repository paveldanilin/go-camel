package camel

import (
	"testing"
)

func TestLoopCountProcessor(t *testing.T) {
	loop := newLoopCountProcessor(5).
		WithStepName("Loop with 5 iterations").
		WithProcessor(newSetBodyProcessor(mustSimpleExpr("exchange.properties.CAMEL_LOOP_INDEX")))

	exchange := NewExchange(nil, nil)

	loop.Process(exchange)

	expectedBody := 4
	if exchange.Message().Body != expectedBody {
		t.Errorf("TestLoopCountProcessor() = %v; want body %v", exchange.Message().Body, expectedBody)
	}
}

func TestLoopWhileProcessor(t *testing.T) {
	loop := newLoopWhileProcessor(mustSimpleExpr("exchange.properties.CAMEL_LOOP_INDEX < 10")).
		WithStepName("Loop with 10 iterations").
		WithProcessor(newSetBodyProcessor(mustSimpleExpr("exchange.properties.CAMEL_LOOP_INDEX")))

	exchange := NewExchange(nil, nil)

	loop.Process(exchange)

	expectedBody := 9
	if exchange.Message().Body != expectedBody {
		t.Errorf("TestLoopWhileProcessor() = %v; want body %v", exchange.Message().Body, expectedBody)
	}
}
