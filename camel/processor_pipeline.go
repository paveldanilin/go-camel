package camel

type pipelineProcessor struct {
	name string

	// If TRUE - exit from pipeline on first error.
	// If FALSE - proceed pipeline when error occurs, thus let process error in the next processor.
	stopOnError bool
	processors  []Processor
}

// Pipeline creates a new pipeline of processors.
// Exit from pipeline in case of any error.
func newPipelineProcessor(name string, stopOnError bool) *pipelineProcessor {
	return &pipelineProcessor{
		name:        name,
		stopOnError: stopOnError,
		processors:  []Processor{},
	}
}

func (p *pipelineProcessor) getName() string {
	return p.name
}

func (p *pipelineProcessor) addProcessor(processor Processor) *pipelineProcessor {
	p.processors = append(p.processors, processor)
	return p
}

func (p *pipelineProcessor) Process(exchange *Exchange) {
	for _, pp := range p.processors {
		pp.Process(exchange)
		if exchange.IsError() && p.stopOnError {
			break
		}
	}
}
