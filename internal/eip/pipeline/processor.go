package pipeline

import (
	"github.com/paveldanilin/go-camel/pkg/camel/api"
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
)

type pipelineProcessor struct {
	routeName string
	name      string
	// If TRUE - exit from pipeline on first error.
	// If FALSE - proceed pipeline when error occurs, thus let process error in the next processor.
	stopOnError bool
	processors  []api.Processor
}

func NewProcessor(routeName, name string, stopOnError bool) *pipelineProcessor {
	return &pipelineProcessor{
		routeName:   routeName,
		name:        name,
		stopOnError: stopOnError,
		processors:  []api.Processor{},
	}
}

func (p *pipelineProcessor) Name() string {
	return p.name
}

func (p *pipelineProcessor) RouteName() string {
	return p.routeName
}

func (p *pipelineProcessor) AddProcessor(processor api.Processor) *pipelineProcessor {
	p.processors = append(p.processors, processor)
	return p
}

func (p *pipelineProcessor) Process(e *exchange.Exchange) {
	for _, pp := range p.processors {
		pp.Process(e)
		if e.IsError() && p.stopOnError {
			break
		}
	}
}
