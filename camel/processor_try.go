package camel

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type errorMatcher func(error) bool

func errorEquals(str string) errorMatcher {
	return func(err error) bool {
		if err == nil {
			return false
		}

		return err.Error() == str
	}
}

func errorContains(str string) errorMatcher {
	substrLower := strings.ToLower(str)

	return func(err error) bool {
		if err == nil {
			return false
		}

		errorMsg := strings.ToLower(err.Error())

		return strings.Contains(errorMsg, substrLower)
	}
}

func errorIs(target string) errorMatcher {
	return func(err error) bool {
		return errors.Is(err, errors.New(target))
	}
}

func errorMatches(pattern string) errorMatcher {
	errRegex := regexp.MustCompile(pattern)

	return func(err error) bool {
		if err == nil {
			return false
		}
		return errRegex.MatchString(err.Error())
	}
}

func errorAny() errorMatcher {
	return func(err error) bool {
		return err != nil
	}
}

type catchClause struct {
	errorMatcher errorMatcher
	handler      Processor
}

type tryProcessor struct {
	id                string
	processors        []Processor
	catchClauses      []catchClause
	finallyProcessors []Processor
}

func newTryProcessor(id string) *tryProcessor {
	return &tryProcessor{
		id:           id,
		processors:   []Processor{},
		catchClauses: []catchClause{},
	}
}

func (p *tryProcessor) getId() string {
	return p.id
}

func (p *tryProcessor) addProcessor(processor Processor) *tryProcessor {
	p.processors = append(p.processors, processor)
	return p
}

func (p *tryProcessor) addCatch(errorMatcher errorMatcher, handler Processor) *tryProcessor {
	p.catchClauses = append(p.catchClauses, catchClause{
		errorMatcher: errorMatcher,
		handler:      handler,
	})
	return p
}

func (p *tryProcessor) addFinally(finally ...Processor) *tryProcessor {
	p.finallyProcessors = append(p.finallyProcessors, finally...)
	return p
}

func (p *tryProcessor) Process(exchange *Exchange) {
	var originalErr error

	// Try-block
	for _, processor := range p.processors {
		if invokeProcessor(processor, exchange) || exchange.IsError() {
			originalErr = exchange.Error()
			break
		}
	}

	// Catch-block
	caught := false
	if originalErr != nil && len(p.catchClauses) > 0 {
		for _, c := range p.catchClauses {
			if c.errorMatcher(originalErr) {
				// Execute handler
				invokeProcessor(c.handler, exchange)
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
			invokeProcessor(p, exchange)
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
