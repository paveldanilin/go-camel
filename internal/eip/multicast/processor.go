package multicast

import (
	"github.com/paveldanilin/go-camel/internal/processor"
	"github.com/paveldanilin/go-camel/pkg/camel/api"
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
	"sync"
)

type multicastProcessor struct {
	routeName   string
	name        string
	parallel    bool
	stopOnError bool
	aggregator  exchange.Aggregator
	outputs     []api.Processor // each output is a start of sub step
}

func NewProcessor(routeName, name string, parallel, stopOnError bool, aggregator exchange.Aggregator) *multicastProcessor {
	return &multicastProcessor{
		routeName:   routeName,
		name:        name,
		parallel:    parallel,
		stopOnError: stopOnError,
		aggregator:  aggregator,
		outputs:     []api.Processor{},
	}
}

func (p *multicastProcessor) Name() string {
	return p.name
}

func (p *multicastProcessor) RouteName() string {
	return p.routeName
}

func (p *multicastProcessor) AddOutput(processor api.Processor) {
	p.outputs = append(p.outputs, processor)
}

func (p *multicastProcessor) Process(e *exchange.Exchange) {
	if p.parallel {
		p.parallelProcess(e)
	} else {
		p.syncProcess(e)
	}
}

func (p *multicastProcessor) syncProcess(e *exchange.Exchange) {
	var oldExchange *exchange.Exchange = nil

	for _, outputProcessor := range p.outputs {
		copyExchange := e.Copy()

		processor.Invoke(outputProcessor, copyExchange)
		if copyExchange.IsError() && p.stopOnError {
			break
		}

		if p.aggregator != nil {
			oldExchange = p.aggregator.AggregateExchange(oldExchange, copyExchange)
		}
	}

	if oldExchange != nil {
		e = oldExchange
	}
}

func (p *multicastProcessor) parallelProcess(e *exchange.Exchange) {
	copyExchanges := make([]*exchange.Exchange, len(p.outputs))
	for i := 0; i < len(p.outputs); i++ {
		copyExchanges[i] = e.Copy()
	}

	// TODO: errgroup?
	var wg sync.WaitGroup
	wg.Add(len(p.outputs))

	for i := 0; i < len(p.outputs); i++ {
		outputProcessor := p.outputs[i]
		ex := copyExchanges[i]

		go func() {
			defer wg.Done()
			processor.Invoke(outputProcessor, ex)
		}()
	}

	wg.Wait()

	if p.aggregator != nil {
		var oldExchange *exchange.Exchange = nil
		for _, ex := range copyExchanges {
			oldExchange = p.aggregator.AggregateExchange(oldExchange, ex)
		}
		if oldExchange != nil {
			e = oldExchange
		}
	}
}
