package test

import (
	"errors"
	"github.com/paveldanilin/go-camel/camel"
	"github.com/paveldanilin/go-camel/camel/expr"
	"github.com/paveldanilin/go-camel/camel/processor"
	"testing"
)

func TestPipelineProcessor(t *testing.T) {

	sum := processor.Pipeline(
		processor.SetHeader("a", expr.Const(1)),
		processor.SetHeader("b", expr.Const(1)),
		processor.Process(func(message *camel.Message) {

			a, _ := message.MessageHeaders().Get("a")
			b, _ := message.MessageHeaders().Get("b")
			message.Body = a.(int) + b.(int)

		}))

	m := camel.NewMessage()

	sum.Process(m)
	if m.IsError() {
		panic(m.Error)
	}

	result := m.Body
	expected := 2

	if result != expected {
		t.Errorf("TestPipelineProcessor() = %d; want %d", result, expected)
	}
}

func TestSetPayloadProcessor(t *testing.T) {

	mul := processor.SetBody(expr.Func(func(message *camel.Message) (any, error) {

		a := message.MustHeader("a")
		b := message.MustHeader("b")

		return a.(int) * b.(int), nil
	}))

	m := camel.NewMessage()
	m.MessageHeaders().Set("a", 2)
	m.MessageHeaders().Set("b", 3)

	mul.Process(m)
	if m.IsError() {
		panic(m.Error)
	}

	result := m.Body
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
			processor.SetHeader("ERROR", expr.MustSimple("properties.error"))).
		Finally(processor.SetBody(expr.Const("RESULT")))

	m := camel.NewMessage()

	tryBlock.Process(m)

	if !errors.Is(m.MustHeader("ERROR").(error), mandatoryParameterMissingErr) {
		t.Errorf("TestDoTryProcessor_Error() = %v; want %v", m.Error, mandatoryParameterMissingErr)
	}

	expected := "RESULT"
	if m.Body != expected {
		t.Errorf("TestDoTryProcessor_Error() = %v; want %s", m.Body, expected)
	}
}

func TestChoiceProcessor(t *testing.T) {

	choice := processor.Choice().
		When(expr.MustSimple("header.val > 5"), processor.SetBody(expr.Const(555))).
		When(expr.MustSimple("header.val < 5"), processor.SetBody(expr.Const(777)))

	m := camel.NewMessage()
	m.SetHeader("val", 2)

	choice.Process(m)

	expected := 777
	if m.Body != expected {
		t.Errorf("TestChoiceProcessor() = %v; want %v", m.Body, expected)
	}
}

func TestLoopProcessor(t *testing.T) {

	loop := processor.Loop(5, nil, processor.LogMessage(">"))

	m := camel.NewMessage()
	loop.Process(m)
}
