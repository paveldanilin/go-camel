package camel

import (
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
)

// RouteStepHook represents pre- or post-processor hook function.
type RouteStepHook func(*exchange.Exchange)

// RouteStep represents a common processor wrapper.
// Holds context attributes like route name and route step for each processor instance.
// Adds common logics for Process func:
//  1. Set Exchange properties: exchange.PropertyRouteName, exchange.PropertyStepName.
//  2. If present exchange.MessageHeaderHistory, adds new record to exchange.MessageHistory.
//  3. Check if context cancelled or timeout reached.
//  4. Add pre/post RouteStepHook.
//
// EACH processor MUST be wrapped by the RouteStep processor.
type RouteStep struct {
	routeName string
	stepName  string
	processor Processor
	preHook   RouteStepHook
	postHook  RouteStepHook
}

func NewRouteStep(routeName, stepName string, processor Processor, preHook RouteStepHook, postHook RouteStepHook) *RouteStep {
	return &RouteStep{
		routeName: routeName,
		stepName:  stepName,
		processor: processor,
		preHook:   preHook,
		postHook:  postHook,
	}
}

func (p *RouteStep) RouteName() string {
	return p.routeName
}

func (p *RouteStep) StepName() string {
	return p.stepName
}

func (p *RouteStep) Process(e *exchange.Exchange) {
	e.SetProperty(exchange.PropertyRouteName, p.RouteName())
	e.SetProperty(exchange.PropertyStepName, p.StepName())

	// Check if exchange supports MessageHistory
	if mh, supportsMessageHistory := e.MessageHeader(exchange.MessageHeaderHistory); supportsMessageHistory {
		if hist, isMessageHistory := mh.(*exchange.MessageHistory); isMessageHistory {
			rec := exchange.NewMessageHistoryRecord(p.RouteName(), p.StepName())
			hist.AddRecord(rec)
			// Update elapsed time
			defer rec.UpdateElapsedTime()
		}
	}

	if p.postHook != nil {
		defer func() {
			p.postHook(e)
		}()
	}

	if err := e.CheckCancelOrTimeout(); err != nil {
		e.SetError(err)
		return
	}

	if p.preHook != nil {
		p.preHook(e)
	}

	p.processor.Process(e)
}
