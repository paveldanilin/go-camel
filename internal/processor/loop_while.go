package processor

import (
	"errors"
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
)

type WhileLoop struct {
	predicate  Predicate
	copy       bool // TRUE - make shallow copy of Exchange for each iteration
	processors []Processor
}

func NewWhileLoop(predicate Predicate, copy bool) (*WhileLoop, error) {
	if predicate == nil {
		return nil, errors.New("predicate is nil")
	}
	return &WhileLoop{
		predicate:  predicate,
		processors: []Processor{},
		copy:       copy,
	}, nil
}

// AddProcessor adds processor to the loop body
func (p *WhileLoop) AddProcessor(processor Processor) *WhileLoop {
	p.processors = append(p.processors, processor)
	return p
}

func (p *WhileLoop) Process(e *exchange.Exchange) {
	if len(p.processors) == 0 {
		return // Nothing to iterate
	}

	iterations := 0
	for {
		e.SetProperty(exchange.PropertyLoopIndex, iterations)

		// User set break loop property
		if v, exists := e.Property(exchange.PropertyLoopBreak); exists && v.(bool) == true {
			break
		}

		// Check while condition
		predicateResult, err := p.predicate.Test(e)
		if err != nil {
			e.SetError(nil)
			break
		}
		if !predicateResult {
			break
		}

		// Make shallow copy
		var currentExchange *exchange.Exchange
		if p.copy {
			cp := *e // Shallow copy
			currentExchange = &cp
		} else {
			currentExchange = e
		}

		// LoopCount through processors
		breakIteration := false
		for _, pp := range p.processors {
			if invoke(pp, e) || currentExchange.Error() != nil {
				breakIteration = true
				break
			}
		}

		// If shallow copy, copy it back
		if p.copy {
			*e = *currentExchange
		}

		// panic/error breaks loop
		if breakIteration || e.IsError() {
			break
		}

		iterations++
	}
}
