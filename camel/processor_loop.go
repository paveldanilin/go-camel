package camel

import (
	"fmt"
)

type loopCountProcessor struct {
	// stepName is a logical name of current operation.
	stepName   string
	count      int
	processors []Processor
	// TRUE - make shallow copy for each iteration
	copy bool
}

func newLoopCountProcessor(count int, processors ...Processor) *loopCountProcessor {
	return &loopCountProcessor{
		stepName:   fmt.Sprintf("loop{count:%d}", count),
		count:      count,
		processors: processors,
		copy:       true,
	}
}

func (p *loopCountProcessor) WithStepName(stepName string) *loopCountProcessor {
	p.stepName = stepName
	return p
}

func (p *loopCountProcessor) WithProcessor(processor Processor) *loopCountProcessor {
	p.processors = append(p.processors, processor)
	return p
}

func (p *loopCountProcessor) Process(exchange *Exchange) {
	if !exchange.On(p.stepName) {
		return
	}

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
			if invokeWithRecovery(processor, exchange) || currentExchange.Error() != nil {
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
	// stepName is a logical name of current operation.
	stepName   string
	predicate  Predicate
	processors []Processor
	// TRUE - make shallow copy for each iteration
	copy bool
}

func newLoopWhileProcessor(predicate Expr, processors ...Processor) *loopWhileProcessor {
	if predicate == nil {
		panic(fmt.Errorf("camel: processor: LoopWhile predicate cannot be nil"))
	}
	return &loopWhileProcessor{
		stepName:   fmt.Sprintf("loop{}"),
		predicate:  newPredicateExpr(predicate),
		processors: processors,
		copy:       true,
	}
}

func (p *loopWhileProcessor) WithStepName(stepName string) *loopWhileProcessor {
	p.stepName = stepName
	return p
}

func (p *loopWhileProcessor) WithProcessor(processor Processor) *loopWhileProcessor {
	p.processors = append(p.processors, processor)
	return p
}

func (p *loopWhileProcessor) Process(exchange *Exchange) {
	if !exchange.On(p.stepName) {
		return
	}

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
			if invokeWithRecovery(processor, exchange) || currentExchange.Error() != nil {
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
