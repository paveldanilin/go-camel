package camel

import (
	"fmt"
	"github.com/paveldanilin/go-camel/pkg/camel/expr"
	"github.com/paveldanilin/go-camel/pkg/camel/routestep"
)

type ChoiceStepBuilder struct {
	builder    *RouteBuilder
	choiceStep *routestep.Choice
}

// When adds ChoiceWhen block to the current ChoiceStep and returns ChoiceStepBuilder.
func (cb *ChoiceStepBuilder) When(predicate expr.Expression, configure func(b *RouteBuilder)) *ChoiceStepBuilder {
	if cb.builder.err != nil {
		return cb
	}

	whenCase := routestep.ChoiceWhen{Predicate: predicate}

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
