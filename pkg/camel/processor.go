package camel

import (
	"fmt"
	"github.com/paveldanilin/go-camel/pkg/camel/api"
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
)

type named interface {
	Name() string
}

func getProcessorName(p api.Processor) string {
	if n, isNamed := p.(named); isNamed {
		return n.Name()
	}
	return fmt.Sprintf("%T", p)
}

type routed interface {
	RouteName() string
}

func getRouteName(p api.Processor) string {
	if r, isRouted := p.(routed); isRouted {
		return r.RouteName()
	}
	return fmt.Sprintf("%T", p)
}

// processor represents a decorator for any processor with pre/post processing functions.
type processor struct {
	delegate      api.Processor
	preProcessor  func(*exchange.Exchange)
	postProcessor func(*exchange.Exchange)
}

func decorateProcessor(p api.Processor, preProcessor func(*exchange.Exchange), postProcessor func(*exchange.Exchange)) *processor {
	return &processor{
		delegate:      p,
		preProcessor:  preProcessor,
		postProcessor: postProcessor,
	}
}

func (p *processor) Process(e *exchange.Exchange) {
	// Check if exchange supports MessageHistory
	if mh, supportsMessageHistory := e.Message().Header(exchange.CamelHeaderMessageHistory); supportsMessageHistory {
		if hist, isMessageHistory := mh.(*exchange.MessageHistory); isMessageHistory {
			rec := exchange.NewMessageHistoryRecord(getRouteName(p.delegate), getProcessorName(p.delegate))
			hist.AddRecord(rec)
			// Update elapsed time
			defer rec.UpdateElapsedTime()
		}
	}

	if err := e.CheckCancelOrTimeout(); err != nil {
		e.SetError(err)
		return
	}

	if p.postProcessor != nil {
		defer func() {
			p.postProcessor(e)
		}()
	}

	if p.preProcessor != nil {
		p.preProcessor(e)
	}

	p.delegate.Process(e)
}
