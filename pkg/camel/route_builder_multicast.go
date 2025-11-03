package camel

import (
	"github.com/paveldanilin/go-camel/pkg/camel/api"
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
	"github.com/paveldanilin/go-camel/pkg/camel/step"
)

type MulticastStepBuilder struct {
	builder       *RouteBuilder
	multicastStep *step.Multicast
}

func (mb *MulticastStepBuilder) ParallelProcessing() *MulticastStepBuilder {
	mb.multicastStep.Parallel = true
	return mb
}

func (mb *MulticastStepBuilder) SyncProcessing() *MulticastStepBuilder {
	mb.multicastStep.Parallel = false
	return mb
}

func (mb *MulticastStepBuilder) StopOnError(stopOnError bool) *MulticastStepBuilder {
	mb.multicastStep.StopOnError = stopOnError
	return mb
}

func (mb *MulticastStepBuilder) Aggregator(aggregator exchange.Aggregator) *MulticastStepBuilder {
	mb.multicastStep.Aggregator = aggregator
	return mb
}

func (mb *MulticastStepBuilder) Process(configure func(b *RouteBuilder)) *MulticastStepBuilder {
	if mb.builder.err != nil {
		return mb
	}

	outputProcess := step.OutputProcess{Steps: []api.RouteStep{}}

	mb.builder.pushStack(&outputProcess.Steps)
	configure(mb.builder)
	mb.builder.popStack()

	mb.multicastStep.Outputs = append(mb.multicastStep.Outputs, outputProcess)

	return mb
}

func (mb *MulticastStepBuilder) EndMulticast() *RouteBuilder {
	return mb.builder
}
