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

type messageHistory struct {
	time        time.Time
	elapsedTime int64
	routeName   string
	stepName    string
}

func (mh *messageHistory) ElapsedTime() int64 {
	return mh.elapsedTime
}

func (mh *messageHistory) Time() time.Time {
	return mh.time
}

func (mh *messageHistory) RouteName() string {
	return ""
}

func (mh *messageHistory) StepName() string {
	return mh.stepName
}

func (mh *messageHistory) Message() *Message {
	return nil
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
	start := time.Now()
	stepName := ""
	if idd, isIdd := p.decoratedProcessor.(identifiable); isIdd {
		stepName = idd.getId()
	}

	mh := &messageHistory{
		time:        start,
		elapsedTime: -1,
		routeName:   "",
		stepName:    stepName,
	}
	exchange.pushMessageHistory(mh)

	defer func() {
		mh.elapsedTime = time.Since(mh.time).Milliseconds()
	}()

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
