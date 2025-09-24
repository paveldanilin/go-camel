package dsl

import "fmt"

type PipelineStep struct {
	Name       string
	StoOnError bool
	Steps      []RouteStep
}

func (s *PipelineStep) StepName() string {
	if s.Name == "" {
		return fmt.Sprintf("pipeline[stopOnError=%v]", s.StoOnError)
	}
	return s.Name
}

// Pipeline adds new pipeline.
// Function configure will be called to configure PipelineStep.
func (b *RouteBuilder) Pipeline(stepName string, stopOnError bool, configure func(b *RouteBuilder)) *RouteBuilder {
	if b.err != nil {
		return b
	}
	pipeline := &PipelineStep{
		Name:       stepName,
		StoOnError: stopOnError,
		Steps:      []RouteStep{},
	}
	b.addStep(pipeline)

	// push pipeline
	b.pushStack(&pipeline.Steps)
	configure(b) // configure PipelineStep
	b.popStack() // pop

	return b
}
