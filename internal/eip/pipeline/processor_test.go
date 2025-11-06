package pipeline

import (
	"github.com/paveldanilin/go-camel/internal/eip/fn"
	"github.com/paveldanilin/go-camel/internal/eip/setheader"
	"github.com/paveldanilin/go-camel/internal/expression"
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
	"testing"
)

func TestPipelineProcessor(t *testing.T) {
	sum := NewProcessor("", "sum", false).
		AddProcessor(setheader.NewProcessor("", "set a", "a", expression.NewConst(1))).
		AddProcessor(setheader.NewProcessor("", "set b", "b", expression.NewConst(1))).
		AddProcessor(fn.NewProcessor("", "calc", func(e *exchange.Exchange) {

			a, _ := e.Message().Header("a")
			b, _ := e.Message().Header("b")
			e.Message().Body = a.(int) + b.(int)

		}))

	e := exchange.NewExchange(nil)

	sum.Process(e)
	if e.IsError() {
		panic(e.Error())
	}

	result := e.Message().Body
	expected := 2

	if result != expected {
		t.Errorf("TestPipelineProcessor() = %d; want %d", result, expected)
	}
}
