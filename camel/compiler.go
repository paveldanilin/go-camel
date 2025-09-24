package camel

import (
	"errors"
	"fmt"
	"github.com/paveldanilin/go-camel/dsl"
	"strings"
)

type compilerConfig struct {
	funcRegistry      FuncRegistry
	preProcessorFunc  func(exchange *Exchange)
	postProcessorFunc func(exchange *Exchange)
}

// compileRoute takes dsl.Route and returns runtime representation of route.
func compileRoute(r *dsl.Route, c compilerConfig) (*route, error) {
	producer, err := compileRouteStep(c, r.Steps...)
	if err != nil {
		return nil, err
	}

	return &route{
		name:     r.Name,
		from:     r.From,
		producer: producer,
	}, nil
}

func compileRouteStep(c compilerConfig, s ...dsl.RouteStep) (Processor, error) {
	if s == nil || len(s) == 0 {
		return nil, errors.New("empty steps")
	}

	if len(s) > 1 {
		pipeline := newPipelineProcessor("", false)
		for _, step := range s {
			p, err := compileRouteStep(c, step)
			if err != nil {
				return nil, err
			}
			pipeline.addProcessor(p)
		}
		return decorateProcessor(pipeline, c.preProcessorFunc, c.postProcessorFunc), nil
	}

	switch t := s[0].(type) {
	case *dsl.SetBodyStep:
		p := newSetBodyProcessor(t.StepName(), compileExpr(t.BodyValue))
		return decorateProcessor(p, c.preProcessorFunc, c.postProcessorFunc), nil
	case *dsl.SetHeaderStep:
		p := newSetHeaderProcessor(t.StepName(), t.HeaderName, compileExpr(t.HeaderValue))
		return decorateProcessor(p, c.preProcessorFunc, c.postProcessorFunc), nil
	case *dsl.SetPropertyStep:
		p := newSetPropertyProcessor(t.StepName(), t.PropertyName, compileExpr(t.PropertyValue))
		return decorateProcessor(p, c.preProcessorFunc, c.postProcessorFunc), nil
	case *dsl.ToStep:
		p := newToProcessor(t.StepName(), t.URI)
		return decorateProcessor(p, c.preProcessorFunc, c.postProcessorFunc), nil
	case *dsl.PipelineStep:
		pipeline := newPipelineProcessor(t.StepName(), t.StoOnError)
		for _, step := range t.Steps {
			p, err := compileRouteStep(c, step)
			if err != nil {
				return nil, err
			}
			pipeline.addProcessor(p)
		}
		return decorateProcessor(pipeline, c.preProcessorFunc, c.postProcessorFunc), nil
	case *dsl.ChoiceStep:
		p := newChoiceProcessor(t.StepName())
		for _, when := range t.WhenCases {
			whenBody, err := compileRouteStep(c, when.Steps...)
			if err != nil {
				return nil, err
			}
			p.addWhen(compileExpr(when.Predicate), whenBody)
		}
		if len(t.Otherwise) > 0 {
			otherwise, err := compileRouteStep(c, t.Otherwise...)
			if err != nil {
				return nil, err
			}
			p.setOtherwise(otherwise)
		}
		return decorateProcessor(p, c.preProcessorFunc, c.postProcessorFunc), nil
	case *dsl.TryStep:
		p := newTryProcessor(t.StepName())

		for _, tryStep := range t.Steps {
			tryProcessor, err := compileRouteStep(c, tryStep)
			if err != nil {
				return nil, err
			}
			p.addProcessor(tryProcessor)
		}

		for _, catch := range t.WhenCatches {
			catchProcessor, err := compileRouteStep(c, catch.Steps...)
			if err != nil {
				return nil, err
			}
			p.addCatch(compileErrMatcher(catch.ErrorMatcher), catchProcessor)
		}

		for _, finallyStep := range t.FinallySteps {
			finallyProcessor, err := compileRouteStep(c, finallyStep)
			if err != nil {
				return nil, err
			}
			p.addFinally(finallyProcessor)
		}
		return decorateProcessor(p, c.preProcessorFunc, c.postProcessorFunc), nil

		// FuncStep compiles dsl.FuncStep -> funcProcessor
	case *dsl.FuncStep:
		if userFunc, isUserFunc := t.Func.(func(*Exchange)); isUserFunc {
			return decorateProcessor(newFuncProcessor(t.StepName(), userFunc), c.preProcessorFunc, c.postProcessorFunc), nil
		}
		if storedFuncName, isStoredFunc := t.Func.(string); isStoredFunc {
			storedFunc := c.funcRegistry.Func(storedFuncName)
			if storedFunc == nil {
				return nil, fmt.Errorf("camel: func step: %s: function not found in registry: %s", t.StepName(), storedFuncName)
			}
			return decorateProcessor(newFuncProcessor(t.StepName(), storedFunc), c.preProcessorFunc, c.postProcessorFunc), nil
		}
		return nil, fmt.Errorf("camel: func step: %s: expected function signature 'func(*Exchange)'", t.StepName())

		// SetErrorStep compiles dsl.SetErrorStep -> setErrorProcessor
	case *dsl.SetErrorStep:
		p := newSetErrorProcessor(t.StepName(), t.Error)
		return decorateProcessor(p, c.preProcessorFunc, c.postProcessorFunc), nil

		// SleepStep compiles dsl.SleepStep -> sleepProcessor
	case *dsl.SleepStep:
		p := newSleepProcessor(t.StepName(), t.Duration)
		return decorateProcessor(p, c.preProcessorFunc, c.postProcessorFunc), nil
	}

	return nil, fmt.Errorf("camel: unknown route step: %T", s[0])
}

func compileExpr(expression dsl.Expression) Expr {
	switch expression.Language {
	case "simple":
		s, err := newSimpleExpr(expression.Expression)
		if err != nil {
			panic(fmt.Errorf("camel: expression: simple: %w", err))
		}
		return s
	case "constant":
		return newConstExpr(expression.Value)
	}

	panic(fmt.Errorf("camel: expression: unknown language: %s", expression.Language))
}

func compileErrMatcher(matcher dsl.ErrorMatcher) errorMatcher {
	if strings.TrimSpace(matcher.Target) == "*" || strings.TrimSpace(matcher.Target) == "" {
		return errorAny()
	}

	switch matcher.MatchMode {
	case dsl.ErrorMatchModeIs:
		return errorIs(matcher.Target)
	case dsl.ErrorMatchModeContains:
		return errorContains(matcher.Target)
	case dsl.ErrorMatchModeEquals:
		return errorEquals(matcher.Target)
	case dsl.ErrorMatchModeRegex:
		return errorMatches(matcher.Target)
	}

	panic(fmt.Errorf("camel: error matcher: unknown mode: %s", matcher.MatchMode))
}
