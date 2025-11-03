package try

import (
	"fmt"
	"github.com/paveldanilin/go-camel/internal/processor"
	"github.com/paveldanilin/go-camel/pkg/camel/api"
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
)

type catchClause struct {
	errorMatcher ErrorMatcher
	handler      api.Processor
}

type tryProcessor struct {
	routeName         string
	name              string
	processors        []api.Processor
	catchClauses      []catchClause
	finallyProcessors []api.Processor
}

func NewProcessor(routeName, name string) *tryProcessor {
	return &tryProcessor{
		routeName:    routeName,
		name:         name,
		processors:   []api.Processor{},
		catchClauses: []catchClause{},
	}
}

func (p *tryProcessor) Name() string {
	return p.name
}

func (p *tryProcessor) RouteName() string {
	return p.routeName
}

func (p *tryProcessor) AddProcessor(processor api.Processor) *tryProcessor {
	p.processors = append(p.processors, processor)
	return p
}

func (p *tryProcessor) AddCatch(errorMatcher ErrorMatcher, handler api.Processor) *tryProcessor {
	p.catchClauses = append(p.catchClauses, catchClause{
		errorMatcher: errorMatcher,
		handler:      handler,
	})
	return p
}

func (p *tryProcessor) AddFinally(finally ...api.Processor) *tryProcessor {
	p.finallyProcessors = append(p.finallyProcessors, finally...)
	return p
}

func (p *tryProcessor) Process(e *exchange.Exchange) {
	var originalErr error

	// Try-block
	for _, pt := range p.processors {
		if processor.Invoke(pt, e) || e.IsError() {
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
				processor.Invoke(c.handler, e)
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
			processor.Invoke(pf, e)
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
