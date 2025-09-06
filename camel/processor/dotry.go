package processor

import (
	"fmt"
	"github.com/paveldanilin/go-camel/camel"
)

type DoTryProcessor struct {
	// stepName is a logical name of current operation.
	stepName          string
	processors        []camel.Processor
	catchClauses      []catchClause
	finallyProcessors []camel.Processor
}

func DoTry(processors ...camel.Processor) *DoTryProcessor {
	return &DoTryProcessor{
		stepName:     "doTry{}",
		processors:   processors,
		catchClauses: []catchClause{},
	}
}

func (p *DoTryProcessor) WithStepName(stepName string) *DoTryProcessor {
	p.stepName = stepName
	return p
}

func (p *DoTryProcessor) Catch(predicate func(error) bool, handler camel.Processor) *DoTryProcessor {
	p.catchClauses = append(p.catchClauses, catchClause{
		predicate: predicate,
		handler:   handler,
	})

	return p
}

func (p *DoTryProcessor) Finally(finally ...camel.Processor) *DoTryProcessor {
	p.finallyProcessors = append(p.finallyProcessors, finally...)
	return p
}

func (p *DoTryProcessor) Process(exchange *camel.Exchange) {
	if !exchange.On(p.stepName) {
		return
	}

	var originalErr error

	// Try-block
	for _, processor := range p.processors {
		if InvokeWithRecovery(processor, exchange) || exchange.IsError() {
			originalErr = exchange.Error()
			break
		}
	}

	// Catch-block
	caught := false
	if originalErr != nil && len(p.catchClauses) > 0 {
		for _, c := range p.catchClauses {
			if c.predicate(originalErr) {
				// Execute handler
				InvokeWithRecovery(c.handler, exchange)
				caught = true
				// Clear error on success handling (Camel-like style).
				exchange.SetError(nil)
				// First match only
				break
			}
		}
	}

	// Finally-block
	if len(p.finallyProcessors) > 0 {
		for _, p := range p.finallyProcessors {
			InvokeWithRecovery(p, exchange)
		}

		// In case of error/panic in finally , combines with originalErr (if any)
		if exchange.IsError() && originalErr != nil && !caught {
			exchange.SetError(fmt.Errorf("original error: %w; finally error: %v", originalErr, exchange.Error))
		}
	}

	// Restore originalErr if catch-block does not catch error
	if originalErr != nil && !caught && exchange.Error() == nil {
		exchange.SetError(originalErr)
	}
}

type catchClause struct {
	predicate func(error) bool
	handler   camel.Processor
}
