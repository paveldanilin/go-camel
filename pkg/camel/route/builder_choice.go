package route

import (
	"fmt"
	"github.com/paveldanilin/go-camel/pkg/camel/route/expr"
	"github.com/paveldanilin/go-camel/pkg/camel/route/step"
)

type ChoiceBuilder struct {
	builder    *Builder
	choiceStep *step.Choice
}

func (cb *ChoiceBuilder) WithName(stepName string) *ChoiceBuilder {
	step.ApplyOptions(cb.choiceStep, step.WithName(stepName))
	return cb
}

// When adds step.ChoiceWhen block to the current choice step and returns ChoiceBuilder.
func (cb *ChoiceBuilder) When(predicate expr.Definition, then func(b *Builder)) *ChoiceBuilder {
	if cb.builder.err != nil {
		return cb
	}

	whenCase := step.ChoiceWhen{Predicate: predicate}

	cb.builder.pushStack(&whenCase.Steps)
	then(cb.builder)
	cb.builder.popStack()

	cb.choiceStep.WhenCases = append(cb.choiceStep.WhenCases, whenCase)

	return cb
}

// Otherwise adds the otherwise (default) block to the current choice step and returns main *Builder.
func (cb *ChoiceBuilder) Otherwise(block func(b *Builder)) *Builder {
	if cb.builder.err != nil {
		return cb.builder
	}
	if cb.choiceStep.Otherwise != nil {
		cb.builder.err = fmt.Errorf("step Choice '%s' already has block Otherwise", cb.choiceStep.Name)
		return cb.builder
	}

	cb.builder.pushStack(&cb.choiceStep.Otherwise)
	block(cb.builder)
	cb.builder.popStack()

	return cb.builder
}

// EndChoice returns the main *Builder.
func (cb *ChoiceBuilder) EndChoice() *Builder {
	return cb.builder
}
