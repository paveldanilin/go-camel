package setheader

import (
	"github.com/paveldanilin/go-camel/internal/expression"
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
	"testing"
)

func TestSetHeaderProcessor(t *testing.T) {
	p := NewProcessor("test", "test", "abc", expression.NewConst(3310))

	e := exchange.NewExchange(nil)

	p.Process(e)

	expectedValue := 3310
	resultValue, exists := e.Message().Header("abc")

	if !exists {
		t.Fatalf("TestSetHeaderProcessor() = missing message header 'abc'")
	}

	if resultValue != expectedValue {
		t.Fatalf("TestSetHeaderProcessor() = %v; want = %v", resultValue, expectedValue)
	}
}
