package loop

import (
	"fmt"
	"github.com/paveldanilin/go-camel/internal/expression"
	"github.com/paveldanilin/go-camel/internal/processor"
	"github.com/paveldanilin/go-camel/pkg/camel/api"
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
)

type countProcessor struct {
	routeName  string
	name       string
	count      int
	copy       bool // TRUE - make shallow copy for each iteration
	processors []api.Processor
}

func NewCountProcessor(routeName, name string, count int) *countProcessor {
	return &countProcessor{
		routeName:  routeName,
		name:       name,
		count:      count,
		copy:       true,
		processors: []api.Processor{},
	}
}

func (p *countProcessor) Name() string {
	return p.name
}

func (p *countProcessor) RouteName() string {
	return p.routeName
}

func (p *countProcessor) AddProcessor(processor api.Processor) *countProcessor {
	p.processors = append(p.processors, processor)
	return p
}

func (p *countProcessor) Process(e *exchange.Exchange) {
	if len(p.processors) == 0 || p.count <= 0 {
		return // Nothing to iterate
	}

	iterations := 0
	for {
		e.SetProperty("CAMEL_LOOP_INDEX", iterations)

		if iterations >= p.count {
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
			if processor.Invoke(pp, e) || currentExchange.Error() != nil {
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

type whileProcessor struct {
	routeName  string
	name       string
	predicate  expression.Predicate
	copy       bool // TRUE - make shallow copy for each iteration
	processors []api.Processor
}

func NewWhileProcessor(routeName, name string, predicate expression.Expression) *whileProcessor {
	if predicate == nil {
		panic(fmt.Errorf("camel: processor: LoopWhile predicate cannot be nil"))
	}
	return &whileProcessor{
		routeName:  routeName,
		name:       name,
		predicate:  expression.NewPredicateFromExpression(predicate),
		processors: []api.Processor{},
		copy:       true,
	}
}

func (p *whileProcessor) Name() string {
	return p.name
}

func (p *whileProcessor) RouteName() string {
	return p.routeName
}

func (p *whileProcessor) AddProcessor(processor api.Processor) *whileProcessor {
	p.processors = append(p.processors, processor)
	return p
}

func (p *whileProcessor) Process(e *exchange.Exchange) {
	if len(p.processors) == 0 {
		return // Nothing to iterate
	}

	iterations := 0
	for {
		e.SetProperty("CAMEL_LOOP_INDEX", iterations)

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
			if processor.Invoke(pp, e) || currentExchange.Error() != nil {
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
