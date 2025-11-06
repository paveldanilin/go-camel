package processor

import (
	"fmt"
	"github.com/paveldanilin/go-camel/pkg/camel/api"
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
)

// Invoke invokes processor with a panic recovery.
// Returns TRUE if panic occurs.
// Returns FALSE if no panic occurs.
func Invoke(p api.Processor, e *exchange.Exchange) (panicked bool) {
	defer func() {
		if r := recover(); r != nil {
			e.SetError(fmt.Errorf("%v", r))
			panicked = true
		}
	}()

	p.Process(e)
	return false
}
