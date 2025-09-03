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
		count:      count,
		processors: processors,
		copy:       true,
	}
}

func (p *LoopCountProcessor) SetStepName(stepName string) *LoopCountProcessor {
	p.stepName = stepName
	return p
}

func (p *LoopCountProcessor) AddProc(processor camel.Processor) *LoopCountProcessor {
	p.processors = append(p.processors, processor)
	return p
}

func (p *LoopCountProcessor) Process(exchange *camel.Exchange) {
	exchange.PushStep(p.stepName)

	if len(p.processors) == 0 || p.count <= 0 {
		return // Nothing to iterate
	}

	if err := exchange.CheckCancelOrTimeout(); err != nil {
		exchange.Error = err
		return
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
			if InvokeWithRecovery(processor, exchange) || currentExchange.Error != nil {
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
		predicate:  expr.Predicate(predicate),
		processors: processors,
		copy:       true,
	}
}

func (p *LoopWhileProcessor) SetStepName(stepName string) *LoopWhileProcessor {
	p.stepName = stepName
	return p
}

func (p *LoopWhileProcessor) AddProc(processor camel.Processor) *LoopWhileProcessor {
	p.processors = append(p.processors, processor)
	return p
}

func (p *LoopWhileProcessor) Process(exchange *camel.Exchange) {
	exchange.PushStep(p.stepName)

	if len(p.processors) == 0 {
		return // Nothing to iterate
	}

	if err := exchange.CheckCancelOrTimeout(); err != nil {
		exchange.Error = err
		return
	}

	iterations := 0
	for {
		exchange.SetProperty("CAMEL_LOOP_INDEX", iterations)

		// Check while condition
		predicateResult, err := p.predicate.Test(exchange)
		if err != nil {
			exchange.Error = err
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
			if InvokeWithRecovery(processor, exchange) || currentExchange.Error != nil {
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
