package test

import (
	"errors"
	"github.com/paveldanilin/go-camel/camel"
	"github.com/paveldanilin/go-camel/camel/expr"
	"github.com/paveldanilin/go-camel/camel/processor"
	"testing"
)

func TestPipelineProcessor(t *testing.T) {
	sum := processor.Pipeline().
		WithStepName("Sum").
		WithProcessor(processor.SetHeader("a", expr.Const(1)).WithStepName("Set 'a' argument")).
		WithProcessor(processor.SetHeader("b", expr.Const(1)).WithStepName("Set 'b' argument")).
		WithProcessor(processor.Process(func(exchange *camel.Exchange) {

			a, _ := exchange.Message().Header("a")
			b, _ := exchange.Message().Header("b")
			exchange.Message().Body = a.(int) + b.(int)

		}))

	exchange := camel.NewExchange(nil, nil)

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

func TestSetPayloadProcessor(t *testing.T) {
	mul := processor.SetBody(expr.Func(func(exchange *camel.Exchange) (any, error) {

		a := exchange.Message().MustHeader("a")
		b := exchange.Message().MustHeader("b")

		return a.(int) * b.(int), nil
	}))

	exchange := camel.NewExchange(nil, nil)
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

func TestDoTryProcessor_Error(t *testing.T) {
	var mandatoryParameterMissingErr = errors.New("mandatory parameter missing")

	tryBlock := processor.DoTry(
		processor.SetHeader("a", expr.Const(1)),
		processor.SetHeader("b", expr.Const(1)),
		processor.SetBody(expr.MustSimple("header.a + header.b")),
		processor.Choice().When(expr.MustSimple("body == 2"), processor.SetError(mandatoryParameterMissingErr)),
	).
		Catch(camel.ErrorContains("mandatory parameter missing"),
			processor.SetHeader("ERROR", expr.MustSimple("exchange.error"))).
		Finally(processor.SetBody(expr.Const("RESULT")))

	exchange := camel.NewExchange(nil, nil)

	tryBlock.Process(exchange)

	if !errors.Is(exchange.Message().MustHeader("ERROR").(error), mandatoryParameterMissingErr) {
		t.Errorf("TestDoTryProcessor_Error() = %v; want %v", exchange.Error, mandatoryParameterMissingErr)
	}

	expected := "RESULT"
	if exchange.Message().Body != expected {
		t.Errorf("TestDoTryProcessor_Error() = %v; want %s", exchange.Message().Body, expected)
	}
}

func TestChoiceProcessor(t *testing.T) {
	choice := processor.Choice().
		When(expr.MustSimple("header.val > 5"), processor.SetBody(expr.Const(555))).
		When(expr.MustSimple("header.val < 5"), processor.SetBody(expr.Const(777)))

	exchange := camel.NewExchange(nil, nil)
	exchange.Message().SetHeader("val", 2)

	choice.Process(exchange)

	expected := 777
	if exchange.Message().Body != expected {
		t.Errorf("TestChoiceProcessor() = %v; want %v", exchange.Message().Body, expected)
	}
}

func TestLoopCountProcessor(t *testing.T) {
	loop := processor.LoopCount(5).
		WithStepName("Loop with 5 iterations").
		WithProcessor(processor.SetBody(expr.MustSimple("exchange.properties.CAMEL_LOOP_INDEX")))

	exchange := camel.NewExchange(nil, nil)

	loop.Process(exchange)

	expectedBody := 4
	if exchange.Message().Body != expectedBody {
		t.Errorf("TestLoopCountProcessor() = %v; want body %v", exchange.Message().Body, expectedBody)
	}
}

func TestLoopWhileProcessor(t *testing.T) {
	loop := processor.LoopWhile(expr.MustSimple("exchange.properties.CAMEL_LOOP_INDEX < 10")).
		WithStepName("Loop with 10 iterations").
		WithProcessor(processor.SetBody(expr.MustSimple("exchange.properties.CAMEL_LOOP_INDEX")))

	exchange := camel.NewExchange(nil, nil)

	loop.Process(exchange)

	expectedBody := 9
	if exchange.Message().Body != expectedBody {
		t.Errorf("TestLoopWhileProcessor() = %v; want body %v", exchange.Message().Body, expectedBody)
	}
}
