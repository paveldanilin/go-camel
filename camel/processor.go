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

type identifiable interface {
	getId() string
}

func getProcessorId(p Processor) string {
	if idd, isIdd := p.(identifiable); isIdd {
		return idd.getId()
	}
	return fmt.Sprintf("%T", p)
}

// processor represents a decorator for any processor with pre/post processing functions.
type processor struct {
	decoratedProcessor Processor
	preProcessorFunc   func(*Exchange)
	postProcessorFunc  func(*Exchange)
}

func decorateProcessor(p Processor, preProcessorFunc func(*Exchange), postProcessorFunc func(*Exchange)) *processor {
	return &processor{
		decoratedProcessor: p,
		preProcessorFunc:   preProcessorFunc,
		postProcessorFunc:  postProcessorFunc,
	}
}

func (p *processor) Process(exchange *Exchange) {
	mh := &MessageHistory{
		time:        time.Now(),
		elapsedTime: -1,
		routeName:   "",
		stepName:    getProcessorId(p.decoratedProcessor),
	}
	exchange.pushMessageHistory(mh)

	defer func() {
		mh.elapsedTime = time.Since(mh.time).Milliseconds()
	}()

	if err := exchange.CheckCancelOrTimeout(); err != nil {
		exchange.SetError(err)
		return
	}

	if p.postProcessorFunc != nil {
		defer func() {
			p.postProcessorFunc(exchange)
		}()
	}

	if p.preProcessorFunc != nil {
		p.preProcessorFunc(exchange)
	}

	p.decoratedProcessor.Process(exchange)
}
