package route

import (
	"fmt"
	"github.com/paveldanilin/go-camel/pkg/camel/route/errmatcher"
	"github.com/paveldanilin/go-camel/pkg/camel/route/step"
)

type TryStepBuilder struct {
	builder *Builder
	tryStep *step.Try
}

// Catch adds 'catch' to the current TryStep.
func (tb *TryStepBuilder) Catch(errorMatcher errmatcher.Definition, configure func(b *Builder)) *TryStepBuilder {
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
func (tb *TryStepBuilder) Finally(configure func(b *Builder)) *Builder {
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

func (tb *TryStepBuilder) EndTry() *Builder {
	return tb.builder
}
