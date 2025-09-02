package processor

import "github.com/paveldanilin/go-camel/camel"

type PipelineConfig struct {
	// If TRUE - exit from pipeline on first error.
	// If FALSE - proceed pipeline when error occurs.
	StopOnError bool
}

type PipelineProcessor struct {
	config     PipelineConfig
	processors []camel.Processor
}

// Pipeline creates a new pipeline of processors.
// Exit from pipeline in case of any error.
func Pipeline(processors ...camel.Processor) *PipelineProcessor {
	return &PipelineProcessor{
		config:     PipelineConfig{StopOnError: true},
		processors: processors,
	}
}

func PipelineWithConfig(config PipelineConfig, processors ...camel.Processor) *PipelineProcessor {
	return &PipelineProcessor{
		config:     config,
		processors: processors,
	}
}

func (p *PipelineProcessor) Process(exchange *camel.Exchange) {
	if err := exchange.CheckCancelOrTimeout(); err != nil {
		exchange.Error = err
		return
	}

	for _, processor := range p.processors {
		processor.Process(exchange)
		if exchange.IsError() && p.config.StopOnError {
			break
		}
	}
}
