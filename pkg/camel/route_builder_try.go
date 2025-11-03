package camel

import (
	"fmt"
	"github.com/paveldanilin/go-camel/pkg/camel/errs"
	"github.com/paveldanilin/go-camel/pkg/camel/step"
)

type TryStepBuilder struct {
	builder *RouteBuilder
	tryStep *step.Try
}

// Catch adds 'catch' to the current TryStep.
func (tb *TryStepBuilder) Catch(errorMatcher errs.Matcher, configure func(b *RouteBuilder)) *TryStepBuilder {
	if tb.builder.err != nil {
		return tb
	}

	catchClause := step.CatchWhen{ErrorMatcher: errorMatcher}

	tb.builder.pushStack(&catchClause.Steps)
	configure(tb.builder)
	tb.builder.popStack()

	tb.tryStep.WhenCatches = append(tb.tryStep.WhenCatches, catchClause)

	return tb // Catch chain
}

// Finally adds 'finally' and returns to the main builder.
func (tb *TryStepBuilder) Finally(configure func(b *RouteBuilder)) *RouteBuilder {
	if tb.builder.err != nil {
		return tb.builder
	}
	if tb.tryStep.FinallySteps != nil {
		tb.builder.err = fmt.Errorf("step Try '%s' already has block Finally", tb.tryStep.Name)
		return tb.builder
	}

	tb.builder.pushStack(&tb.tryStep.FinallySteps)
	configure(tb.builder)
	tb.builder.popStack()

	return tb.builder // main builder
}

func (tb *TryStepBuilder) EndTry() *RouteBuilder {
	return tb.builder
}
