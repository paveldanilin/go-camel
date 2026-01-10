package processor

import (
	"github.com/paveldanilin/go-camel/pkg/camel/api"
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
)

type Pipeline struct {
	// If TRUE - exit from pipeline on first error.
	// If FALSE - proceed pipeline when error occurs, thus let process error in the next processor.
	stopOnError bool
	processors  []api.Processor
}

func NewPipeline(stopOnError bool) *Pipeline {
	return &Pipeline{
		stopOnError: stopOnError,
		processors:  []api.Processor{},
	}
}

func (p *Pipeline) AddProcessor(processor api.Processor) *Pipeline {
	p.processors = append(p.processors, processor)
	return p
}

func (p *Pipeline) Process(e *exchange.Exchange) {
	for _, pp := range p.processors {
		pp.Process(e)
		if e.IsError() && p.stopOnError {
			break
		}
	}
}
