package processor

import (
	"fmt"
	"github.com/paveldanilin/go-camel/camel"
)

type DoTryProcessor struct {
	processors        []camel.Processor
	catchClauses      []catchClause
	finallyProcessors []camel.Processor
}

func DoTry(processors ...camel.Processor) *DoTryProcessor {

	return &DoTryProcessor{
		processors:   processors,
		catchClauses: []catchClause{},
	}
}

type catchClause struct {
	predicate func(error) bool
	handler   camel.Processor
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

func (p *DoTryProcessor) Process(message *camel.Message) {

	var originalErr error

	// Try-block
	for _, processor := range p.processors {
		if Invoke(processor, message) || message.IsError() {
			originalErr = message.Error
			break
		}
	}

	// Catch-block
	caught := false
	if originalErr != nil && len(p.catchClauses) > 0 {
		for _, c := range p.catchClauses {
			if c.predicate(originalErr) {
				// Execute handler
				Invoke(c.handler, message)
				caught = true
				// Clear error on success handling (Camel-like style).
				message.Error = nil
				// First match only
				break
			}
		}
	}

	// Finally-block
	if len(p.finallyProcessors) > 0 {
		for _, p := range p.finallyProcessors {
			Invoke(p, message)
		}

		// In case of error/panic in finally , combines with originalErr (if any)
		if message.IsError() && originalErr != nil && !caught {
			message.Error = fmt.Errorf("original error: %w; finally error: %v", originalErr, message.Error)
		}
	}

	// Restore originalErr if catch-block does not catch error
	if originalErr != nil && !caught && message.Error == nil {
		message.Error = originalErr
	}
}
