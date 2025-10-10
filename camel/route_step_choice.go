package camel

import "fmt"

type ChoiceWhen struct {
	Predicate Expression
	Steps     []RouteStep
}

func (s *ChoiceWhen) StepName() string {
	return fmt.Sprintf("choiceWhen[%s:%s]", s.Predicate.Language, s.Predicate.Expression)
}

type ChoiceStep struct {
	Name      string
	WhenCases []ChoiceWhen
	Otherwise []RouteStep
}

func (s *ChoiceStep) StepName() string { return s.Name }

// ---------------------------------------------------------------------------------------------------------------------
// ChoiceStepBuilder
// ---------------------------------------------------------------------------------------------------------------------

type ChoiceStepBuilder struct {
	builder    *RouteBuilder
	choiceStep *ChoiceStep
}

// When adds ChoiceWhen block to the current ChoiceStep and returns ChoiceStepBuilder.
func (cb *ChoiceStepBuilder) When(predicate Expression, configure func(b *RouteBuilder)) *ChoiceStepBuilder {
	if cb.builder.err != nil {
		return cb
	}

	whenCase := ChoiceWhen{Predicate: predicate}

	cb.builder.pushStack(&whenCase.Steps)
	configure(cb.builder)
	cb.builder.popStack()

	cb.choiceStep.WhenCases = append(cb.choiceStep.WhenCases, whenCase)

	return cb
}

// Otherwise adds the otherwise (default) block to the current ChoiceStep and returns main RouteBuilder.
func (cb *ChoiceStepBuilder) Otherwise(configure func(b *RouteBuilder)) *RouteBuilder {
	if cb.builder.err != nil {
		return cb.builder
	}
	if cb.choiceStep.Otherwise != nil {
		cb.builder.err = fmt.Errorf("step Choice '%s' already has block Otherwise", cb.choiceStep.Name)
		return cb.builder
	}

	cb.builder.pushStack(&cb.choiceStep.Otherwise)
	configure(cb.builder)
	cb.builder.popStack()

	return cb.builder
}

// EndChoice returns the main RouteBuilder.
func (cb *ChoiceStepBuilder) EndChoice() *RouteBuilder {
	return cb.builder
}

// ---------------------------------------------------------------------------------------------------------------------
// RouteBuilder :: Choice
// ---------------------------------------------------------------------------------------------------------------------

// Choice adds choice step to the current route level.
func (b *RouteBuilder) Choice(stepName string) *ChoiceStepBuilder {
	if b.err != nil {
		return &ChoiceStepBuilder{builder: b}
	}

	choiceStep := &ChoiceStep{Name: stepName}
	b.addStep(choiceStep)

	return &ChoiceStepBuilder{builder: b, choiceStep: choiceStep}
}
