package processor

import (
	"github.com/paveldanilin/go-camel/internal/expression"
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
	"testing"
)

func TestSetPropertyProcessor(t *testing.T) {
	p := NewSetProperty("MY_PROPERTY", expression.NewConstant("MY_PROPERTY_VALUE"))

	e := exchange.NewExchange(nil)

	p.Process(e)

	resultValue, _ := e.Property("MY_PROPERTY")
	expectedValue := "MY_PROPERTY_VALUE"

	if resultValue != expectedValue {
		t.Fatalf("TestSetPropertyProcessor() = %v; want = %v", resultValue, expectedValue)
	}
}
