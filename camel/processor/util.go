package processor

import (
	"fmt"
	"github.com/paveldanilin/go-camel/camel"
)

// InvokeWithRecovery invokes processor with panic recovery.
// Returns TRUE if panic occurs.
// Returns FALSE if no panic occurs.
func InvokeWithRecovery(p camel.Processor, exchange *camel.Exchange) (panicked bool) {
	defer func() {
		if r := recover(); r != nil {
			exchange.Error = fmt.Errorf("panic recovered: %v", r)
			panicked = true
		}
	}()

	p.Process(exchange)
	return false
}
