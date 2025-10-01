package dsl

import "fmt"

type ChoiceWhen struct {
	Predicate Expression
	Steps     []RouteStep
}

func (s *ChoiceWhen) StepName() string {
	return fmt.Sprintf("when[%s:%s]", s.Predicate.Language, s.Predicate.Expression)
}

type ChoiceStep struct {
	Name      string
	WhenCases []ChoiceWhen
	Otherwise []RouteStep
}

func (s *ChoiceStep) StepName() string { return s.Name }

type ChoiceStepBuilder struct {
	builder    *RouteBuilder
	choiceStep *ChoiceStep
}

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

	return cb.builder // main builder
}

func (cb *ChoiceStepBuilder) EndChoice() *RouteBuilder {
	return cb.builder
}

/*
// When добавляет условную ветку в последний созданный ChoiceWhen.
func (b *RouteBuilder) When(predicate Expression, configure func(b *RouteBuilder)) *RouteBuilder {
	if b.err != nil {
		return b
	}
	if b.lastChoice == nil {
		b.err = fmt.Errorf("bad method call: When() must be called only inside Choice()")
		return b
	}

	whenCase := ChoiceWhen{Predicate: predicate, Steps: []RouteStep{}}
	b.lastChoice.WhenCases = append(b.lastChoice.WhenCases, whenCase)

	stepsPtr := &b.lastChoice.WhenCases[len(b.lastChoice.WhenCases)-1].Steps

	b.pushStack(stepsPtr)
	configure(b) // configure ChoiceWhen
	b.popStack()

	return b
}

func (b *RouteBuilder) Otherwise(configure func(b *RouteBuilder)) *RouteBuilder {
	if b.err != nil {
		return b
	}
	if b.lastChoice == nil {
		b.err = fmt.Errorf("bad method call: Otherwise() must be called only inside ChoiceStepKind()")
		return b
	}
	if b.lastChoice.Otherwise != nil {
		b.err = fmt.Errorf("could not redefine Otherwise since it already created for ChoiceStepKind")
		return b
	}

	b.lastChoice.Otherwise = []RouteStep{}
	b.pushStack(&b.lastChoice.Otherwise)
	configure(b) // configure Otherwise
	b.popStack()

	return b
}
*/

// Choice adds choice step.
// Function configure will be called to configure ChoiceStep.
func (b *RouteBuilder) Choice(stepName string) *ChoiceStepBuilder {
	if b.err != nil {
		return &ChoiceStepBuilder{builder: b}
	}

	choiceStep := &ChoiceStep{Name: stepName}
	b.addStep(choiceStep)

	//b.pushStack(&choiceStep.WhenCases)
	//configure(b) // configure choice
	//b.popStack()

	return &ChoiceStepBuilder{builder: b, choiceStep: choiceStep}

	//choice := &ChoiceStep{Name: stepName}
	//b.addStep(choice)

	// store link to the current choice
	//prevChoice := b.lastChoice
	//b.lastChoice = choice

	//configure(b) // configure ChoiceWhen/Otherwise

	// restore link to the external choice (if any)
	//b.lastChoice = prevChoice
	//return b
}
