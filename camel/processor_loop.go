package camel

import (
	"fmt"
)

type loopCountProcessor struct {
	id string

	count      int
	copy       bool // TRUE - make shallow copy for each iteration
	processors []Processor
}

func newLoopCountProcessor(id string, count int) *loopCountProcessor {
	return &loopCountProcessor{
		id:         id,
		count:      count,
		copy:       true,
		processors: []Processor{},
	}
}

func (p *loopCountProcessor) addProcessor(processor Processor) *loopCountProcessor {
	p.processors = append(p.processors, processor)
	return p
}

func (p *loopCountProcessor) Process(exchange *Exchange) {
	if len(p.processors) == 0 || p.count <= 0 {
		return // Nothing to iterate
	}

	iterations := 0
	for {
		exchange.SetProperty("CAMEL_LOOP_INDEX", iterations)

		if iterations >= p.count {
			break
		}

		// Make shallow copy
		var currentExchange *Exchange
		if p.copy {
			cp := *exchange // Shallow copy
			currentExchange = &cp
		} else {
			currentExchange = exchange
		}

		// LoopCount through processors
		breakIteration := false
		for _, pp := range p.processors {
			if invokeProcessor(pp, exchange) || currentExchange.Error() != nil {
				breakIteration = true
				break
			}
		}

		// If shallow copy, copy it back
		if p.copy {
			*exchange = *currentExchange
		}

		// panic/error breaks loop
		if breakIteration || exchange.IsError() {
			break
		}

		iterations++
	}
}

type loopWhileProcessor struct {
	id string

	predicate  Predicate
	copy       bool // TRUE - make shallow copy for each iteration
	processors []Processor
}

func newLoopWhileProcessor(id string, predicate expression) *loopWhileProcessor {
	if predicate == nil {
		panic(fmt.Errorf("camel: processor: LoopWhile predicate cannot be nil"))
	}
	return &loopWhileProcessor{
		id:         id,
		predicate:  newPredicateFromExpression(predicate),
		processors: []Processor{},
		copy:       true,
	}
}

func (p *loopWhileProcessor) addProcessor(processor Processor) *loopWhileProcessor {
	p.processors = append(p.processors, processor)
	return p
}

func (p *loopWhileProcessor) Process(exchange *Exchange) {
	if len(p.processors) == 0 {
		return // Nothing to iterate
	}

	iterations := 0
	for {
		exchange.SetProperty("CAMEL_LOOP_INDEX", iterations)

		// Check while condition
		predicateResult, err := p.predicate.Test(exchange)
		if err != nil {
			exchange.SetError(nil)
			break
		}
		if !predicateResult {
			break
		}

		// Make shallow copy
		var currentExchange *Exchange
		if p.copy {
			cp := *exchange // Shallow copy
			currentExchange = &cp
		} else {
			currentExchange = exchange
		}

		// LoopCount through processors
		breakIteration := false
		for _, pp := range p.processors {
			if invokeProcessor(pp, exchange) || currentExchange.Error() != nil {
				breakIteration = true
				break
			}
		}

		// If shallow copy, copy it back
		if p.copy {
			*exchange = *currentExchange
		}

		// panic/error breaks loop
		if breakIteration || exchange.IsError() {
			break
		}

		iterations++
	}
}
