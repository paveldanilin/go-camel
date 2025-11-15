package processor

import (
	"fmt"
	"github.com/paveldanilin/go-camel/pkg/camel/api"
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
)

type named interface {
	Name() string
}

type routed interface {
	RouteName() string
}

type Processor struct {
	delegate      api.Processor
	preProcessor  func(*exchange.Exchange)
	postProcessor func(*exchange.Exchange)
}

func NewProcessor(p api.Processor, preProcessor func(*exchange.Exchange), postProcessor func(*exchange.Exchange)) *Processor {
	return &Processor{
		delegate:      p,
		preProcessor:  preProcessor,
		postProcessor: postProcessor,
	}
}

func (p *Processor) Name() string {
	if n, isNamed := p.delegate.(named); isNamed {
		return n.Name()
	}
	return fmt.Sprintf("%T", p.delegate)
}

func (p *Processor) RouteName() string {
	if r, isRouted := p.delegate.(routed); isRouted {
		return r.RouteName()
	}
	return fmt.Sprintf("%T", p.delegate)
}

func (p *Processor) Process(e *exchange.Exchange) {
	// Check if exchange supports MessageHistory
	if mh, supportsMessageHistory := e.Message().Header(exchange.CamelHeaderMessageHistory); supportsMessageHistory {
		if hist, isMessageHistory := mh.(*exchange.MessageHistory); isMessageHistory {
			rec := exchange.NewMessageHistoryRecord(p.RouteName(), p.Name())
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
