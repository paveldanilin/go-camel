package camel

type pipelineProcessor struct {
	id string

	// If TRUE - exit from pipeline on first error.
	// If FALSE - proceed pipeline when error occurs, thus let process error in the next processor.
	stopOnError bool
	processors  []Processor
}

// Pipeline creates a new pipeline of processors.
// Exit from pipeline in case of any error.
func newPipelineProcessor(id string, stopOnError bool) *pipelineProcessor {
	return &pipelineProcessor{
		id:          id,
		stopOnError: stopOnError,
		processors:  []Processor{},
	}
}

func (p *pipelineProcessor) getId() string {
	return p.id
}

func (p *pipelineProcessor) addProcessor(processor Processor) *pipelineProcessor {
	p.processors = append(p.processors, processor)
	return p
}

func (p *pipelineProcessor) Process(exchange *Exchange) {
	for _, processor := range p.processors {
		processor.Process(exchange)
		if exchange.IsError() && p.stopOnError {
			break
		}
	}
}
