package camel

import (
	"fmt"
)

type loopCountProcessor struct {
	id string

	count      int
	processors []Processor
	copy       bool // TRUE - make shallow copy for each iteration
}

func newLoopCountProcessor(id string, count int) *loopCountProcessor {
	return &loopCountProcessor{
		id:         id,
		count:      count,
		processors: []Processor{},
		copy:       true,
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
			copy := *exchange // Shallow copy
			currentExchange = &copy
		} else {
			currentExchange = exchange
		}

		// LoopCount through processors
		breakIteration := false
		for _, processor := range p.processors {
			if invokeProcessor(processor, exchange) || currentExchange.Error() != nil {
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
	processors []Processor
	copy       bool // TRUE - make shallow copy for each iteration
}

func newLoopWhileProcessor(id string, predicate Expr) *loopWhileProcessor {
	if predicate == nil {
		panic(fmt.Errorf("camel: processor: LoopWhile predicate cannot be nil"))
	}
	return &loopWhileProcessor{
		id:         id,
		predicate:  newPredicateFromExpr(predicate),
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
			copy := *exchange // Shallow copy
			currentExchange = &copy
		} else {
			currentExchange = exchange
		}

		// LoopCount through processors
		breakIteration := false
		for _, processor := range p.processors {
			if invokeProcessor(processor, exchange) || currentExchange.Error() != nil {
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
