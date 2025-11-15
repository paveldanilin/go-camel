package route

import (
	"github.com/paveldanilin/go-camel/pkg/camel/api"
	"github.com/paveldanilin/go-camel/pkg/camel/route/step"
)

type MulticastStepBuilder struct {
	builder       *Builder
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

func (mb *MulticastStepBuilder) Aggregator(aggregator api.ExchangeAggregator) *MulticastStepBuilder {
	mb.multicastStep.Aggregator = aggregator
	return mb
}

func (mb *MulticastStepBuilder) Process(configure func(b *Builder)) *MulticastStepBuilder {
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

func (mb *MulticastStepBuilder) EndMulticast() *Builder {
	return mb.builder
}
