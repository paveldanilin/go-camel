package delay

import (
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
	"testing"
	"time"
)

func TestDelayProcessor(t *testing.T) {
	p := NewProcessor("test", "test", 500)
	e := exchange.NewExchange(nil)
	start := time.Now()

	p.Process(e)

	elapsedMs := time.Since(start).Milliseconds()

	expectedValue := int64(500)
	if elapsedMs < expectedValue {
		t.Fatalf("TestDelayProcessor() = %d elapsed ms; want >= %d", elapsedMs, expectedValue)
	}
}
