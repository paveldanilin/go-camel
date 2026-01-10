package processor

import (
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
	"testing"
)

func TestFnProcessor(t *testing.T) {
	p := NewFn(func(e *exchange.Exchange) {
		e.Message().SetBody("Hello, World!")
	})

	e := exchange.NewExchange(nil)
	p.Process(e)

	expectedValue := "Hello, World!"
	if e.Message().Body() != expectedValue {
		t.Fatalf("TestFnProcessor() = %v, want = %s", e.Message().Body(), expectedValue)
	}
}
