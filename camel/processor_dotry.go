package camel

import (
	"fmt"
)

type catchClause struct {
	predicate func(error) bool
	handler   Processor
}

type doTryProcessor struct {
	// stepName is a logical name of current operation.
	stepName          string
	processors        []Processor
	catchClauses      []catchClause
	finallyProcessors []Processor
}

func newDoTryProcessor(processors ...Processor) *doTryProcessor {
	return &doTryProcessor{
		stepName:     "doTry{}",
		processors:   processors,
		catchClauses: []catchClause{},
	}
}

func (p *doTryProcessor) WithStepName(stepName string) *doTryProcessor {
	p.stepName = stepName
	return p
}

func (p *doTryProcessor) Catch(predicate func(error) bool, handler Processor) *doTryProcessor {
	p.catchClauses = append(p.catchClauses, catchClause{
		predicate: predicate,
		handler:   handler,
	})

	return p
}

func (p *doTryProcessor) Finally(finally ...Processor) *doTryProcessor {
	p.finallyProcessors = append(p.finallyProcessors, finally...)
	return p
}

func (p *doTryProcessor) Process(exchange *Exchange) {
	if !exchange.On(p.stepName) {
		return
	}

	var originalErr error

	// Try-block
	for _, processor := range p.processors {
		if invokeWithRecovery(processor, exchange) || exchange.IsError() {
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
				invokeWithRecovery(c.handler, exchange)
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
			invokeWithRecovery(p, exchange)
		}

		// In case of error/panic in finally , combines with originalErr (if any)
		if exchange.IsError() && originalErr != nil && !caught {
			exchange.SetError(fmt.Errorf("original error: %w; finally error: %v", originalErr, exchange.Error()))
		}
	}

	// Restore originalErr if catch-block does not catch error
	if originalErr != nil && !caught && exchange.Error() == nil {
		exchange.SetError(originalErr)
	}
}
