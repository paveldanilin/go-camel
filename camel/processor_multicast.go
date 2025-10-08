package camel

import "sync"

type ExchangeAggregator interface {
	AggregateExchange(oldExchange *Exchange, newExchange *Exchange) *Exchange
}

type multicastProcessor struct {
	name        string
	parallel    bool
	stopOnError bool
	aggregator  ExchangeAggregator
	outputs     []Processor // each output is a start of sub route
}

func newMulticastProcessor(name string, parallel, stopOnError bool, aggregator ExchangeAggregator) *multicastProcessor {
	return &multicastProcessor{
		name:        name,
		parallel:    parallel,
		stopOnError: stopOnError,
		aggregator:  aggregator,
		outputs:     []Processor{},
	}
}

func (p *multicastProcessor) getName() string {
	return p.name
}

func (p *multicastProcessor) addOutput(processor Processor) {
	p.outputs = append(p.outputs, processor)
}

func (p *multicastProcessor) Process(exchange *Exchange) {
	if p.parallel {
		p.parallelProcess(exchange)
	} else {
		p.syncProcess(exchange)
	}
}

func (p *multicastProcessor) syncProcess(exchange *Exchange) {
	var oldExchange *Exchange = nil

	for _, outputProcessor := range p.outputs {
		copyExchange := exchange.Copy()

		invokeProcessor(outputProcessor, copyExchange)
		if copyExchange.IsError() && p.stopOnError {
			break
		}

		if p.aggregator != nil {
			oldExchange = p.aggregator.AggregateExchange(oldExchange, copyExchange)
		}
	}

	if oldExchange != nil {
		exchange = oldExchange
	}
}

func (p *multicastProcessor) parallelProcess(exchange *Exchange) {
	copyExchanges := make([]*Exchange, len(p.outputs))
	for i := 0; i < len(p.outputs); i++ {
		copyExchanges[i] = exchange.Copy()
	}

	// TODO: errgroup?
	var wg sync.WaitGroup
	wg.Add(len(p.outputs))

	for i := 0; i < len(p.outputs); i++ {
		outputProcessor := p.outputs[i]
		ex := copyExchanges[i]

		go func() {
			defer wg.Done()
			invokeProcessor(outputProcessor, ex)
		}()
	}

	wg.Wait()

	if p.aggregator != nil {
		var oldExchange *Exchange = nil
		for _, ex := range copyExchanges {
			oldExchange = p.aggregator.AggregateExchange(oldExchange, ex)
		}
		if oldExchange != nil {
			exchange = oldExchange
		}
	}
}
