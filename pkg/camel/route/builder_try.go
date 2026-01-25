package route

import (
	"fmt"
	"github.com/paveldanilin/go-camel/pkg/camel/route/errmatcher"
	"github.com/paveldanilin/go-camel/pkg/camel/route/step"
)

type TryBuilder struct {
	builder *Builder
	tryStep *step.Try
}

func (tb *TryBuilder) WithName(stepName string) *TryBuilder {
	step.ApplyOptions(tb.tryStep, step.WithName(stepName))
	return tb
}

// Catch adds 'catch' to the current TryStep.
func (tb *TryBuilder) Catch(errorMatcher errmatcher.Definition, then StepBuilder) *TryBuilder {
	if tb.builder.err != nil {
		return tb
	}

	catchClause := step.CatchWhen{ErrorMatcher: errorMatcher}

	tb.builder.pushStack(&catchClause.Steps)
	then(tb.builder)
	tb.builder.popStack()

	tb.tryStep.WhenCatches = append(tb.tryStep.WhenCatches, catchClause)

	return tb // Catch chain
}

func (tb *TryBuilder) CatchAny(then StepBuilder) *TryBuilder {
	return tb.Catch(errmatcher.Any(), then)
}

func (tb *TryBuilder) CatchErrorIs(target error, then StepBuilder) *TryBuilder {
	return tb.Catch(errmatcher.Is(target), then)
}

func (tb *TryBuilder) CatchErrorEquals(str string, then StepBuilder) *TryBuilder {
	return tb.Catch(errmatcher.Equals(str), then)
}

func (tb *TryBuilder) CatchErrorContains(str string, then StepBuilder) *TryBuilder {
	return tb.Catch(errmatcher.Contains(str), then)
}

func (tb *TryBuilder) CatchErrorMatches(pattern string, then StepBuilder) *TryBuilder {
	return tb.Catch(errmatcher.Matches(pattern), then)
}

// Finally adds 'finally' and returns to the main builder.
func (tb *TryBuilder) Finally(block StepBuilder) *Builder {
	if tb.builder.err != nil {
		return tb.builder
	}
	if tb.tryStep.FinallySteps != nil {
		tb.builder.err = fmt.Errorf("step Try '%s' already has block Finally", tb.tryStep.Name)
		return tb.builder
	}

	tb.builder.pushStack(&tb.tryStep.FinallySteps)
	block(tb.builder)
	tb.builder.popStack()

	return tb.builder // main builder
}

func (tb *TryBuilder) EndTry() *Builder {
	return tb.builder
}
