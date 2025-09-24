package dsl

import (
	"fmt"
	"strings"
)

type CatchWhen struct {
	ErrorMatcher ErrorMatcher
	Steps        []RouteStep
}

type TryStep struct {
	Name         string
	Steps        []RouteStep
	WhenCatches  []CatchWhen
	FinallySteps []RouteStep
}

func (s *TryStep) StepName() string {
	if s.Name == "" {
		when := make([]string, len(s.WhenCatches))
		for i, w := range s.WhenCatches {
			when[i] = fmt.Sprintf("%s:%s", w.ErrorMatcher.MatchMode, w.ErrorMatcher.Target)
		}
		return fmt.Sprintf("try[%s]", strings.Join(when, ";"))
	}
	return s.Name
}

type TryStepBuilder struct {
	builder *RouteBuilder
	tryStep *TryStep
}

// Catch adds 'catch' to the current TryStep.
func (tb *TryStepBuilder) Catch(errorMatcher ErrorMatcher, configure func(b *RouteBuilder)) *TryStepBuilder {
	if tb.builder.err != nil {
		return tb
	}

	catchClause := CatchWhen{ErrorMatcher: errorMatcher}

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

// Try adds try-catch-finally step.
func (b *RouteBuilder) Try(stepName string, configure func(b *RouteBuilder)) *TryStepBuilder {
	if b.err != nil {
		return &TryStepBuilder{builder: b}
	}

	tryStep := &TryStep{Name: stepName}
	b.addStep(tryStep)

	b.pushStack(&tryStep.Steps)
	configure(b) // configure try
	b.popStack()

	return &TryStepBuilder{builder: b, tryStep: tryStep}
}

type ErrorMatchMode string

const (
	ErrorMatchModeEquals   ErrorMatchMode = "equals"
	ErrorMatchModeContains ErrorMatchMode = "contains"
	ErrorMatchModeRegex    ErrorMatchMode = "regex"
	ErrorMatchModeIs       ErrorMatchMode = "is"
)

type ErrorMatcher struct {
	MatchMode ErrorMatchMode
	Target    string
}

func ErrIs(target error) ErrorMatcher {
	return ErrorMatcher{
		MatchMode: ErrorMatchModeIs,
		Target:    target.Error(),
	}
}

func ErrEquals(str string) ErrorMatcher {
	return ErrorMatcher{
		MatchMode: ErrorMatchModeEquals,
		Target:    str,
	}
}

func ErrAny() ErrorMatcher {
	return ErrorMatcher{
		MatchMode: ErrorMatchModeEquals,
		Target:    "*",
	}
}

func ErrContains(str string) ErrorMatcher {
	return ErrorMatcher{
		MatchMode: ErrorMatchModeContains,
		Target:    str,
	}
}

func ErrMatches(pattern string) ErrorMatcher {
	return ErrorMatcher{
		MatchMode: ErrorMatchModeRegex,
		Target:    pattern,
	}
}
