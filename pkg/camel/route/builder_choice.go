package route

import (
	"fmt"
	"github.com/paveldanilin/go-camel/pkg/camel/route/expr"
	"github.com/paveldanilin/go-camel/pkg/camel/route/step"
)

type ChoiceStepBuilder struct {
	builder    *Builder
	choiceStep *step.Choice
}

// When adds ChoiceWhen block to the current ChoiceStep and returns ChoiceStepBuilder.
func (cb *ChoiceStepBuilder) When(predicate expr.Definition, configure func(b *Builder)) *ChoiceStepBuilder {
	if cb.builder.err != nil {
		return cb
	}

	whenCase := step.ChoiceWhen{Predicate: predicate}

	cb.builder.pushStack(&whenCase.Steps)
	configure(cb.builder)
	cb.builder.popStack()

	cb.choiceStep.WhenCases = append(cb.choiceStep.WhenCases, whenCase)

	return cb
}

// Otherwise adds the otherwise (default) block to the current ChoiceStep and returns main RouteBuilder.
func (cb *ChoiceStepBuilder) Otherwise(configure func(b *Builder)) *Builder {
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
func (cb *ChoiceStepBuilder) EndChoice() *Builder {
	return cb.builder
}
