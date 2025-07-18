package processor

import (
	"fmt"
	"github.com/paveldanilin/go-camel/camel"
)

type DoTryProcessor struct {
	tryProcessor     camel.Processor
	catchClauses     []catchClause
	finallyProcessor camel.Processor
}

func DoTry(try camel.Processor) *DoTryProcessor {

	return &DoTryProcessor{
		tryProcessor: try,
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

func (p *DoTryProcessor) Finally(finally camel.Processor) *DoTryProcessor {

	p.finallyProcessor = finally

	return p
}

func (p *DoTryProcessor) Process(message *camel.Message) {

	// Catches panic
	safeProcess := func(p camel.Processor) (panicked bool) {

		defer func() {
			if r := recover(); r != nil {
				message.SetError(fmt.Errorf("panic recovered: %v", r))
				panicked = true
			}
		}()

		p.Process(message)
		return false
	}

	// Original error
	var originalErr error

	// Try-block
	if safeProcess(p.tryProcessor) {
		originalErr = message.Error()
	}

	// Catch-block
	caught := false
	if originalErr != nil && len(p.catchClauses) > 0 {
		for _, c := range p.catchClauses {
			if c.predicate(originalErr) {
				// Execute handler
				safeProcess(c.handler)
				caught = true
				// Clear error on success handling (Camel-like style).
				message.SetError(nil)
				// First match only
				break
			}
		}
	}

	// Finally-block
	if p.finallyProcessor != nil {
		safeProcess(p.finallyProcessor) // Collect error/panic

		// In case of error/panic in finally , combines with originalErr (if any)
		if message.IsError() && originalErr != nil && !caught {
			message.SetError(fmt.Errorf("original error: %w; finally error: %v", originalErr, message.Error()))
		}
	}

	// Restore originalErr if catch-block does not catch error
	if originalErr != nil && !caught && message.Error() == nil {
		message.SetError(originalErr)
	}
}
