package fn

import (
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
	"testing"
)

func TestFnProcessor(t *testing.T) {
	p := NewProcessor("", "", func(e *exchange.Exchange) {
		e.Message().Body = "Hello, World!"
	})

	e := exchange.NewExchange(nil)
	p.Process(e)

	expectedValue := "Hello, World!"
	if e.Message().Body != expectedValue {
		t.Fatalf("TestFnProcessor() = %v, want = %s", e.Message().Body, expectedValue)
	}
}
