package processor

import (
	"fmt"
	"github.com/paveldanilin/go-camel/pkg/camel/api"
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
)

type CatchClause struct {
	errorMatcher ErrorMatcher
	handler      api.Processor
}

type Try struct {
	processors        []api.Processor
	catchClauses      []CatchClause
	finallyProcessors []api.Processor
}

func NewTry() *Try {
	return &Try{
		processors:   []api.Processor{},
		catchClauses: []CatchClause{},
	}
}

func (p *Try) AddProcessor(processor api.Processor) *Try {
	p.processors = append(p.processors, processor)
	return p
}

func (p *Try) AddCatch(errorMatcher ErrorMatcher, handler api.Processor) *Try {
	p.catchClauses = append(p.catchClauses, CatchClause{
		errorMatcher: errorMatcher,
		handler:      handler,
	})
	return p
}

func (p *Try) AddFinally(finally ...api.Processor) *Try {
	p.finallyProcessors = append(p.finallyProcessors, finally...)
	return p
}

func (p *Try) Process(e *exchange.Exchange) {
	var originalErr error

	// Try-block
	for _, pt := range p.processors {
		if invoke(pt, e) || e.IsError() {
			originalErr = e.Error()
			break
		}
	}

	// Catch-block
	caught := false
	if originalErr != nil && len(p.catchClauses) > 0 {
		for _, c := range p.catchClauses {
			if c.errorMatcher(originalErr) {
				// Execute handler
				invoke(c.handler, e)
				caught = true
				// Clear error on success handling (Camel-like style).
				e.SetError(nil)
				// First match only
				break
			}
		}
	}

	// Finally-block
	if len(p.finallyProcessors) > 0 {
		for _, pf := range p.finallyProcessors {
			invoke(pf, e)
		}

		// In case of error/panic in finally , combines with originalErr (if any)
		if e.IsError() && originalErr != nil && !caught {
			e.SetError(fmt.Errorf("original error: %w; finally error: %v", originalErr, e.Error()))
		}
	}

	// Restore originalErr if catch-block does not catch error
	if originalErr != nil && !caught && e.Error() == nil {
		e.SetError(originalErr)
	}
}
