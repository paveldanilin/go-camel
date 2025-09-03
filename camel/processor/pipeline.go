package processor

import "github.com/paveldanilin/go-camel/camel"

type PipelineProcessor struct {
	// stepName is a logical name of current operation.
	stepName string
	// If TRUE - exit from pipeline on first error.
	// If FALSE - proceed pipeline when error occurs, thus let process error in the next processor.
	stopOnError bool
	processors  []camel.Processor
}

// Pipeline creates a new pipeline of processors.
// Exit from pipeline in case of any error.
func Pipeline(processors ...camel.Processor) *PipelineProcessor {
	return &PipelineProcessor{
		stopOnError: true,
		processors:  processors,
	}
}

func (p *PipelineProcessor) SetStepName(stepName string) *PipelineProcessor {
	p.stepName = stepName
	return p
}

func (p *PipelineProcessor) SetStopOnError(stopOnError bool) *PipelineProcessor {
	p.stopOnError = stopOnError
	return p
}

func (p *PipelineProcessor) AddProc(processor camel.Processor) *PipelineProcessor {
	p.processors = append(p.processors, processor)
	return p
}

func (p *PipelineProcessor) Process(exchange *camel.Exchange) {
	exchange.PushStep(p.stepName)

	if err := exchange.CheckCancelOrTimeout(); err != nil {
		exchange.Error = err
		return
	}

	for _, processor := range p.processors {
		processor.Process(exchange)
		if exchange.IsError() && p.stopOnError {
			break
		}
	}
}
