package processor

import (
	"github.com/paveldanilin/go-camel/internal/expression"
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
	"testing"
)

func TestLoopCountProcessor(t *testing.T) {
	loop := NewRepeatLoop(5).
		AddProcessor(NewSetBody(expression.MustSimple("property.CAMEL_LOOP_INDEX")))

	e := exchange.NewExchange(nil)

	loop.Process(e)

	expectedBody := 4
	if e.Message().Body() != expectedBody {
		t.Errorf("TestLoopCountProcessor() = %v; want body %v", e.Message().Body(), expectedBody)
	}
}
