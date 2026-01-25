package processor

import (
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
)

type RepeatLoop struct {
	count      int
	copy       bool // TRUE - make shallow copy for each iteration
	processors []Processor
}

func NewRepeatLoop(count int) *RepeatLoop {
	return &RepeatLoop{
		count:      count,
		copy:       true,
		processors: []Processor{},
	}
}

// AddProcessor adds processor to the loop body
func (p *RepeatLoop) AddProcessor(processor Processor) *RepeatLoop {
	p.processors = append(p.processors, processor)
	return p
}

func (p *RepeatLoop) Process(e *exchange.Exchange) {
	if len(p.processors) == 0 || p.count <= 0 {
		return // Nothing to iterate
	}

	iterations := 0
	for {
		e.SetProperty(exchange.PropertyLoopIndex, iterations)

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
