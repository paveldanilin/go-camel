package processor

import (
	"github.com/paveldanilin/go-camel/internal/expression"
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
	"testing"
)

func TestWhileLoopProcessor(t *testing.T) {
	loop, err := NewWhileLoop(expression.MustSimple("property.CAMEL_LOOP_INDEX < 10"), true)
	if err != nil {
		t.Fatalf("Failed to create while loop processor: %s", err)
	}

	loop.AddProcessor(NewSetBody(expression.MustSimple("property.CAMEL_LOOP_INDEX")))

	e := exchange.NewExchange(nil)

	loop.Process(e)

	expectedBody := 9
	if e.Message().Body() != expectedBody {
		t.Errorf("TestWhileLoopProcessor() = %v; want body %v", e.Message().Body(), expectedBody)
	}
}

func TestWhileLoopProcessor_Break(t *testing.T) {
	loop, err := NewWhileLoop(expression.MustSimple("true"), true)
	if err != nil {
		t.Fatalf("Failed to create while loop processor: %s", err)
	}

	loop.AddProcessor(NewSetBody(expression.MustSimple("property.CAMEL_LOOP_INDEX"))).
		AddProcessor(
			NewChoice().
				// break after 5 iterations
				When(
					expression.MustSimple("body >= 5"),
					NewSetProperty(exchange.PropertyLoopBreak, expression.NewConstant(true))),
		)

	e := exchange.NewExchange(nil)

	loop.Process(e)

	expectedBody := 5
	if e.Message().Body() != expectedBody {
		t.Errorf("TestWhileLoopProcessor_Break() = %v; want body %v", e.Message().Body(), expectedBody)
	}
}
