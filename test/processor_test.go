package test

import (
	"github.com/paveldanilin/go-camel/camel"
	"github.com/paveldanilin/go-camel/camel/expr"
	"github.com/paveldanilin/go-camel/camel/processor"
	"testing"
)

func TestPipelineSum(t *testing.T) {

	sum := processor.Pipeline(
		processor.SetHeader("a", expr.Const(1)),
		processor.SetHeader("b", expr.Const(1)),
		processor.Process(func(message *camel.Message) error {

			a, _ := message.Header("a")
			b, _ := message.Header("b")
			message.SetPayload(a.(int) + b.(int))

			return nil
		}))

	m := camel.NewMessage()

	err := sum.Process(m)
	if err != nil {
		panic(err)
	}

	result := m.Payload()
	expected := 2

	if result != expected {
		t.Errorf("PipelineSum() = %d; want %d", result, expected)
	}
}

func TestSetBodyMul(t *testing.T) {

	mul := processor.SetPayload(expr.Func(func(message *camel.Message) (any, error) {

		a, _ := message.Header("a")
		b, _ := message.Header("b")

		return a.(int) * b.(int), nil
	}))

	m := camel.NewMessage()
	m.SetHeader("a", 2)
	m.SetHeader("b", 3)

	err := mul.Process(m)
	if err != nil {
		panic(err)
	}

	result := m.Payload()
	expected := 6

	if result != expected {
		t.Errorf("SetBodyMul() = %d; want %d", result, expected)
	}
}
