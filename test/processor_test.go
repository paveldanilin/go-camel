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
			message.SetPayload(a.(int) + b.(int))

		}))

	m := camel.NewMessage()

	sum.Process(m)
	if m.IsError() {
		panic(m.Error())
	}

	result := m.Payload()
	expected := 2

	if result != expected {
		t.Errorf("TestPipelineProcessor() = %d; want %d", result, expected)
	}
}

func TestSetBodyProcessor(t *testing.T) {

	mul := processor.SetPayload(expr.Func(func(message *camel.Message) (any, error) {

		a := message.MustHeader("a")
		b := message.MustHeader("b")

		return a.(int) * b.(int), nil
	}))

	m := camel.NewMessage()
	m.MessageHeaders().Set("a", 2)
	m.MessageHeaders().Set("b", 3)

	mul.Process(m)
	if m.IsError() {
		panic(m.Error())
	}

	result := m.Payload()
	expected := 6

	if result != expected {
		t.Errorf("TestSetBodyProcessor() = %d; want %d", result, expected)
	}
}

func TestDoTryProcessor_Error(t *testing.T) {

	var mandatoryParameterMissingErr = errors.New("mandatory parameter missing")

	p := processor.DoTry(processor.Pipeline(
		processor.SetHeader("a", expr.Const(1)),
		processor.SetHeader("b", expr.Const(1)),
		processor.SetPayload(expr.Func(func(message *camel.Message) (any, error) {
			return message.MustHeader("a").(int) + message.MustHeader("b").(int), nil
		})),
		processor.Process(func(message *camel.Message) {
			if !message.HasHeader("factor") {
				message.SetError(mandatoryParameterMissingErr)
			}
		}),
		processor.Process(func(message *camel.Message) {
			result := message.MustHeader("factor").(int) * message.Payload().(int)
			message.SetPayload(result)
		}),
	)).
		Catch(camel.PredicateErrorContains("mandatory parameter missing"), processor.Process(func(message *camel.Message) {
			message.SetHeader("ERROR", message.Error())
		})).
		Finally(processor.Process(func(message *camel.Message) {
			message.SetPayload("RESULT")
		}))

	m := camel.NewMessage()

	p.Process(m)

	expected := "RESULT"
	if m.Payload() != expected {
		t.Errorf("TestDoTryProcessor_Error() = %v; want %s", m.Payload(), expected)
	}
}
