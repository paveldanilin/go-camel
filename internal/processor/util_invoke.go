package processor

import (
	"fmt"
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
)

// invoke processor with a panic recovery.
// Returns TRUE if panic occurs.
// Returns FALSE if no panic occurs.
func invoke(p Processor, e *exchange.Exchange) (panicked bool) {
	defer func() {
		if r := recover(); r != nil {
			e.SetError(fmt.Errorf("%v", r))
			panicked = true
		}
	}()

	p.Process(e)
	return false
}
