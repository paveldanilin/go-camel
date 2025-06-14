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

func (p *PipelineProcessor) Process(message *camel.Message) error {

	for _, processor := range p.processors {
		err := processor.Process(message)
		if err != nil {
			return err
		}
	}

	return nil
}
