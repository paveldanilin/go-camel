package processor

import (
	"github.com/paveldanilin/go-camel/pkg/camel/api"
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
)

// Hook represents pre- or post-processor hook function.
type Hook func(exchange2 *exchange.Exchange)

// RouteStep represents a common processor wrapper.
// Holds context attributes like route name and route step for each processor instance.
// Adds common logics for Process func:
//  1. Set Exchange properties: exchange.PropertyRouteName, exchange.PropertyStepName.
//  2. If present exchange.MessageHeaderHistory, adds new record to exchange.MessageHistory.
//  3. Check if context cancelled or timeout reached.
//  4. Add pre/post Hook.
//
// EACH processor MUST be wrapped by the RouteStep processor.
type RouteStep struct {
	routeName string
	stepName  string
	proc      api.Processor
	preHook   Hook
	postHook  Hook
}

func NewRouteStep(routeName, stepName string, proc api.Processor, preHook Hook, postHook Hook) *RouteStep {
	return &RouteStep{
		routeName: routeName,
		stepName:  stepName,
		proc:      proc,
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
	if mh, supportsMessageHistory := e.Message().Header(exchange.MessageHeaderHistory); supportsMessageHistory {
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

	p.proc.Process(e)
}
