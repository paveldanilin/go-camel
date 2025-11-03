package loop

import (
	"github.com/paveldanilin/go-camel/internal/eip/setbody"
	"github.com/paveldanilin/go-camel/internal/expression"
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
	"testing"
)

func TestLoopCountProcessor(t *testing.T) {
	loop := NewCountProcessor("", "Loop with 5 iterations", 5).
		AddProcessor(setbody.NewProcessor("", "set body", expression.MustSimple("property.CAMEL_LOOP_INDEX")))

	e := exchange.NewExchange(nil)

	loop.Process(e)

	expectedBody := 4
	if e.Message().Body != expectedBody {
		t.Errorf("TestLoopCountProcessor() = %v; want body %v", e.Message().Body, expectedBody)
	}
}

func TestLoopWhileProcessor(t *testing.T) {
	loop := NewWhileProcessor("", "Loop with 10 iterations", expression.MustSimple("property.CAMEL_LOOP_INDEX < 10")).
		AddProcessor(setbody.NewProcessor("", "set body", expression.MustSimple("property.CAMEL_LOOP_INDEX")))

	e := exchange.NewExchange(nil)

	loop.Process(e)

	expectedBody := 9
	if e.Message().Body != expectedBody {
		t.Errorf("TestLoopWhileProcessor() = %v; want body %v", e.Message().Body, expectedBody)
	}
}
