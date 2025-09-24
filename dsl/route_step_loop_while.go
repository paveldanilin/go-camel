package dsl

import "fmt"

type LoopWhileStep struct {
	Name         string
	Predicate    Expression
	CopyExchange bool
	Steps        []RouteStep
}

func (s *LoopWhileStep) StepName() string {
	if s.Name == "" {
		return fmt.Sprintf("loop[%s:%s]", s.Predicate.Language, s.Predicate.Expression)
	}
	return s.Name
}

func (b *RouteBuilder) LoopWhile(stepName string, predicate Expression, copyExchange bool, configure func(b *RouteBuilder)) *RouteBuilder {
	if b.err != nil {
		return b
	}

	step := &LoopWhileStep{
		Name:         stepName,
		Predicate:    predicate,
		CopyExchange: copyExchange,
	}
	b.addStep(step)

	b.pushStack(&step.Steps)
	configure(b)
	b.popStack()

	return b
}
