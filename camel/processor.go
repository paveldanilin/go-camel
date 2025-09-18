package camel

import "fmt"

// invokeWithRecovery invokes processor with panic recovery.
// Returns TRUE if panic occurs.
// Returns FALSE if no panic occurs.
func invokeWithRecovery(p Processor, exchange *Exchange) (panicked bool) {
	defer func() {
		if r := recover(); r != nil {
			exchange.SetError(fmt.Errorf("panic recovered: %v", r))
			panicked = true
		}
	}()

	p.Process(exchange)
	return false
}

// processor represents a wrapper for any processor with pre/post processing functions.
type processor struct {
	Processor         Processor
	preProcessorFunc  func(*Exchange)
	postProcessorFunc func(*Exchange)
}

func newProcessor(p Processor, preProcessorFunc func(*Exchange), postProcessorFunc func(*Exchange)) *processor {
	return &processor{
		Processor:         p,
		preProcessorFunc:  preProcessorFunc,
		postProcessorFunc: postProcessorFunc,
	}
}

func (p processor) Process(exchange *Exchange) {
	if p.postProcessorFunc != nil {
		defer func() {
			p.postProcessorFunc(exchange)
		}()
	}

	if p.preProcessorFunc != nil {
		p.preProcessorFunc(exchange)
	}

	p.Processor.Process(exchange)
}
