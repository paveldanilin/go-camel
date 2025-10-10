package camel

import (
	"fmt"
	"time"
)

// invokeProcessor invokes processor with a panic recovery.
// Returns TRUE if panic occurs.
// Returns FALSE if no panic occurs.
func invokeProcessor(p Processor, exchange *Exchange) (panicked bool) {
	defer func() {
		if r := recover(); r != nil {
			exchange.SetError(fmt.Errorf("%v", r))
			panicked = true
		}
	}()

	p.Process(exchange)
	return false
}

type named interface {
	getName() string
}

func getProcessorName(p Processor) string {
	if n, isNamed := p.(named); isNamed {
		return n.getName()
	}
	return fmt.Sprintf("%T", p)
}

// processor represents a decorator for any processor with pre/post processing functions.
type processor struct {
	decoratedProcessor Processor
	preProcessor       func(*Exchange)
	postProcessor      func(*Exchange)
}

func decorateProcessor(p Processor, preProcessor func(*Exchange), postProcessor func(*Exchange)) *processor {
	return &processor{
		decoratedProcessor: p,
		preProcessor:       preProcessor,
		postProcessor:      postProcessor,
	}
}

func (p *processor) Process(exchange *Exchange) {
	mh := &MessageHistory{
		time:        time.Now(),
		elapsedTime: -1,
		routeName:   "",
		stepName:    getProcessorName(p.decoratedProcessor),
	}
	exchange.pushMessageHistory(mh)

	defer func() {
		mh.elapsedTime = time.Since(mh.time).Milliseconds()
	}()

	if err := exchange.CheckCancelOrTimeout(); err != nil {
		exchange.SetError(err)
		return
	}

	if p.postProcessor != nil {
		defer func() {
			p.postProcessor(exchange)
		}()
	}

	if p.preProcessor != nil {
		p.preProcessor(exchange)
	}

	p.decoratedProcessor.Process(exchange)
}
