package camel

import (
	"fmt"
	"strings"
)

type MulticastOutput struct {
	Steps []RouteStep
}

func (sr MulticastOutput) StepName() string {
	return "output"
}

type MulticastStep struct {
	Name        string
	Parallel    bool
	StopOnError bool
	Aggregator  ExchangeAggregator
	Outputs     []MulticastOutput
}

func (s *MulticastStep) StepName() string {
	if s.Name == "" {
		outputNames := make([]string, len(s.Outputs))
		for i, output := range s.Outputs {
			if len(output.Steps) > 0 {
				outputNames[i] = output.Steps[0].StepName()
			} else {
				outputNames[i] = fmt.Sprintf("%d", i)
			}
		}
		return fmt.Sprintf("multicast[parallel=%v;stopOnError=%v;outputs=%s]",
			s.Parallel, s.StopOnError, strings.Join(outputNames, ";"))
	}
	return s.Name
}

// ---------------------------------------------------------------------------------------------------------------------
// MulticastStepBuilder
// ---------------------------------------------------------------------------------------------------------------------

type MulticastStepBuilder struct {
	builder       *RouteBuilder
	multicastStep *MulticastStep
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

func (mb *MulticastStepBuilder) Aggregator(aggregator ExchangeAggregator) *MulticastStepBuilder {
	mb.multicastStep.Aggregator = aggregator
	return mb
}

func (mb *MulticastStepBuilder) Output(configure func(b *RouteBuilder)) *MulticastStepBuilder {
	if mb.builder.err != nil {
		return mb
	}

	subRoute := MulticastOutput{Steps: []RouteStep{}}

	mb.builder.pushStack(&subRoute.Steps)
	configure(mb.builder)
	mb.builder.popStack()

	mb.multicastStep.Outputs = append(mb.multicastStep.Outputs, subRoute)

	return mb
}

func (mb *MulticastStepBuilder) EndMulticast() *RouteBuilder {
	return mb.builder
}

// ---------------------------------------------------------------------------------------------------------------------
// RouteBuilder :: Multicast
// ---------------------------------------------------------------------------------------------------------------------

func (b *RouteBuilder) Multicast(stepName string) *MulticastStepBuilder {
	if b.err != nil {
		return &MulticastStepBuilder{builder: b}
	}

	multicastStep := &MulticastStep{
		Name:        stepName,
		Parallel:    false,
		StopOnError: false,
	}
	b.addStep(multicastStep)

	return &MulticastStepBuilder{builder: b, multicastStep: multicastStep}
}
