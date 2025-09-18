package camel

type pipelineProcessor struct {
	// stepName is a logical name of current operation.
	stepName string
	// If TRUE - exit from pipeline on first error.
	// If FALSE - proceed pipeline when error occurs, thus let process error in the next processor.
	stopOnError bool
	processors  []Processor
}

// Pipeline creates a new pipeline of processors.
// Exit from pipeline in case of any error.
func newPipelineProcessor(processors ...Processor) *pipelineProcessor {
	return &pipelineProcessor{
		stepName:    "pipeline{}",
		stopOnError: true,
		processors:  processors,
	}
}

func (p *pipelineProcessor) WithStepName(stepName string) *pipelineProcessor {
	p.stepName = stepName
	return p
}

func (p *pipelineProcessor) WithStopOnError(stopOnError bool) *pipelineProcessor {
	p.stopOnError = stopOnError
	return p
}

func (p *pipelineProcessor) WithProcessor(processor Processor) *pipelineProcessor {
	p.processors = append(p.processors, processor)
	return p
}

func (p *pipelineProcessor) Process(exchange *Exchange) {
	if !exchange.On(p.stepName) {
		return
	}

	for _, processor := range p.processors {
		processor.Process(exchange)
		if exchange.IsError() && p.stopOnError {
			break
		}
	}
}
