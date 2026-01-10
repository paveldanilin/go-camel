package processor

import (
	"github.com/paveldanilin/go-camel/pkg/camel/api"
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
	"sync"
)

type Multicast struct {
	parallel    bool
	stopOnError bool
	aggregator  api.ExchangeAggregator
	outputs     []api.Processor // each output is a start of sub route
}

func NewMulticast(parallel, stopOnError bool, aggregator api.ExchangeAggregator) *Multicast {
	return &Multicast{
		parallel:    parallel,
		stopOnError: stopOnError,
		aggregator:  aggregator,
		outputs:     []api.Processor{},
	}
}

func (p *Multicast) AddOutput(processor api.Processor) {
	p.outputs = append(p.outputs, processor)
}

func (p *Multicast) Process(e *exchange.Exchange) {
	if p.parallel {
		p.parallelProcess(e)
	} else {
		p.syncProcess(e)
	}
}

func (p *Multicast) syncProcess(e *exchange.Exchange) {
	var oldExchange *exchange.Exchange = nil

	for _, outputProcessor := range p.outputs {
		copyExchange := e.Copy()

		invoke(outputProcessor, copyExchange)
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

func (p *Multicast) parallelProcess(e *exchange.Exchange) {
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
			invoke(outputProcessor, ex)
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
