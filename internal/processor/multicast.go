package processor

import (
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
	"sync"
)

type ExchangeAggregator interface {
	AggregateExchange(oldExchange *exchange.Exchange, newExchange *exchange.Exchange) *exchange.Exchange
}

type Multicast struct {
	parallel    bool
	stopOnError bool
	aggregator  ExchangeAggregator
	processors  []Processor // each output is a start of sub route
}

func NewMulticast(parallel, stopOnError bool, aggregator ExchangeAggregator) *Multicast {
	return &Multicast{
		parallel:    parallel,
		stopOnError: stopOnError,
		aggregator:  aggregator,
		processors:  []Processor{},
	}
}

// AddProcessor adds new output chain
func (p *Multicast) AddProcessor(processor Processor) {
	p.processors = append(p.processors, processor)
}

func (p *Multicast) Process(e *exchange.Exchange) {
	var result *exchange.Exchange
	if p.parallel {
		result = p.parallelProcess(e)
	} else {
		result = p.syncProcess(e)
	}

	// If address has changed, copy it
	if result != nil && result != e {
		*e = *result
	}
}

func (p *Multicast) syncProcess(e *exchange.Exchange) *exchange.Exchange {
	var oldExchange *exchange.Exchange = nil

	for _, outputProcessor := range p.processors {
		// Make copy of the original Exchange for each output processor
		copyExchange := e.Copy()

		// Invoke processor in safe panic-recovery mode
		invoke(outputProcessor, copyExchange)

		if copyExchange.IsError() && p.stopOnError {
			// Stop processing on first error
			return copyExchange
		}

		// Apply aggregation strategy if defined
		if p.aggregator != nil {
			oldExchange = p.aggregator.AggregateExchange(oldExchange, copyExchange)
			if oldExchange != nil && oldExchange.IsError() && p.stopOnError {
				// Stop processing on first error
				return oldExchange
			}
		}
	}

	if oldExchange == nil {
		// By default, we will get the same Exchange instance
		return e
	}

	return oldExchange
}

func (p *Multicast) parallelProcess(e *exchange.Exchange) *exchange.Exchange {
	// Make shallow copy for all outputs
	copyExchanges := make([]*exchange.Exchange, len(p.processors))
	for i := 0; i < len(p.processors); i++ {
		copyExchanges[i] = e.Copy()
	}

	// TODO: errgroup?
	var wg sync.WaitGroup
	wg.Add(len(p.processors))

	for i := 0; i < len(p.processors); i++ {
		outputProcessor := p.processors[i]
		ex := copyExchanges[i]

		go func() {
			defer wg.Done()
			invoke(outputProcessor, ex)
		}()
	}

	wg.Wait()

	// Apply aggregation strategy if defined
	if p.aggregator != nil {
		var oldExchange *exchange.Exchange = nil
		for _, ex := range copyExchanges {
			oldExchange = p.aggregator.AggregateExchange(oldExchange, ex)
		}
		return oldExchange
	}

	// By default, we will get the same Exchange instance
	return e
}
