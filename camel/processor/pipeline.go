package processor

import "github.com/paveldanilin/go-camel/camel"

type PipelineProcessor struct {
	processors []camel.Processor
}

func Pipeline(processors ...camel.Processor) *PipelineProcessor {
	return &PipelineProcessor{
		processors: processors,
	}
}

func (p *PipelineProcessor) Process(message *camel.Message) {

	for _, processor := range p.processors {
		processor.Process(message)
		if message.IsError() {
			break
		}
	}
}
