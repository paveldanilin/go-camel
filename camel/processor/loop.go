package processor

import (
	"fmt"
	"github.com/paveldanilin/go-camel/camel"
	"github.com/paveldanilin/go-camel/camel/expr"
)

type LoopCountProcessor struct {
	// stepName is a logical name of current operation.
	stepName   string
	count      int
	processors []camel.Processor
	// TRUE - make shallow copy for each iteration
	copy bool
}

func LoopCount(count int, processors ...camel.Processor) *LoopCountProcessor {
	return &LoopCountProcessor{
		stepName:   fmt.Sprintf("loop{count:%d}", count),
		count:      count,
		processors: processors,
		copy:       true,
	}
}

func (p *LoopCountProcessor) WithStepName(stepName string) *LoopCountProcessor {
	p.stepName = stepName
	return p
}

func (p *LoopCountProcessor) WithProcessor(processor camel.Processor) *LoopCountProcessor {
	p.processors = append(p.processors, processor)
	return p
}

func (p *LoopCountProcessor) Process(exchange *camel.Exchange) {
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
		var currentExchange *camel.Exchange
		if p.copy {
			copy := *exchange // Shallow copy
			currentExchange = &copy
		} else {
			currentExchange = exchange
		}

		// LoopCount through processors
		breakIteration := false
		for _, processor := range p.processors {
			if InvokeWithRecovery(processor, exchange) || currentExchange.Error() != nil {
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

type LoopWhileProcessor struct {
	// stepName is a logical name of current operation.
	stepName   string
	predicate  camel.Predicate
	processors []camel.Processor
	// TRUE - make shallow copy for each iteration
	copy bool
}

func LoopWhile(predicate camel.Expr, processors ...camel.Processor) *LoopWhileProcessor {
	if predicate == nil {
		panic(fmt.Errorf("camel: processor: LoopWhile predicate cannot be nil"))
	}
	return &LoopWhileProcessor{
		stepName:   fmt.Sprintf("loop{}"),
		predicate:  expr.Predicate(predicate),
		processors: processors,
		copy:       true,
	}
}

func (p *LoopWhileProcessor) WithStepName(stepName string) *LoopWhileProcessor {
	p.stepName = stepName
	return p
}

func (p *LoopWhileProcessor) WithProcessor(processor camel.Processor) *LoopWhileProcessor {
	p.processors = append(p.processors, processor)
	return p
}

func (p *LoopWhileProcessor) Process(exchange *camel.Exchange) {
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
		var currentExchange *camel.Exchange
		if p.copy {
			copy := *exchange // Shallow copy
			currentExchange = &copy
		} else {
			currentExchange = exchange
		}

		// LoopCount through processors
		breakIteration := false
		for _, processor := range p.processors {
			if InvokeWithRecovery(processor, exchange) || currentExchange.Error() != nil {
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
