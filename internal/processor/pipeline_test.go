package processor

import (
	"github.com/paveldanilin/go-camel/internal/expression"
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
	"testing"
)

func TestPipelineProcessor(t *testing.T) {
	sum := NewPipeline(false).
		AddProcessor(NewSetHeader("a", expression.NewConstant(1))).
		AddProcessor(NewSetHeader("b", expression.NewConstant(1))).
		AddProcessor(NewFn(func(e *exchange.Exchange) {

			a, _ := e.Message().Header("a")
			b, _ := e.Message().Header("b")
			e.Message().SetBody(a.(int) + b.(int))

		}))

	e := exchange.NewExchange(nil)

	sum.Process(e)
	if e.IsError() {
		panic(e.Error())
	}

	result := e.Message().Body()
	expected := 2

	if result != expected {
		t.Errorf("TestPipelineProcessor() = %d; want %d", result, expected)
	}
}
