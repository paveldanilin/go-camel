package route

import (
	"github.com/paveldanilin/go-camel/pkg/camel/route/step"
)

type MulticastBuilder struct {
	builder       *Builder
	multicastStep *step.Multicast
}

func (mb *MulticastBuilder) WithName(stepName string) *MulticastBuilder {
	step.ApplyOptions(mb.multicastStep, step.WithName(stepName))
	return mb
}

func (mb *MulticastBuilder) WithParallelProcessing() *MulticastBuilder {
	mb.multicastStep.Parallel = true
	return mb
}

func (mb *MulticastBuilder) WithSyncProcessing() *MulticastBuilder {
	mb.multicastStep.Parallel = false
	return mb
}

func (mb *MulticastBuilder) WithStopOnError(stopOnError bool) *MulticastBuilder {
	mb.multicastStep.StopOnError = stopOnError
	return mb
}

func (mb *MulticastBuilder) WithAggregator(aggregator step.ExchangeAggregator) *MulticastBuilder {
	mb.multicastStep.Aggregator = aggregator
	return mb
}

func (mb *MulticastBuilder) Process(block StepBuilder) *MulticastBuilder {
	if mb.builder.err != nil {
		return mb
	}

	outputProcess := step.OutputProcess{Steps: []step.RouteStep{}}

	mb.builder.pushStack(&outputProcess.Steps)
	block(mb.builder)
	mb.builder.popStack()

	mb.multicastStep.Outputs = append(mb.multicastStep.Outputs, outputProcess)

	return mb
}

func (mb *MulticastBuilder) EndMulticast() *Builder {
	return mb.builder
}
