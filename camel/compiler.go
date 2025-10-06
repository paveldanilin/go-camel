package camel

import (
	"errors"
	"fmt"
	"strings"
)

type compilerConfig struct {
	funcRegistry       FuncRegistry
	dataFormatRegistry DataFormatRegistry
	logger             Logger
	preProcessorFunc   func(exchange *Exchange)
	postProcessorFunc  func(exchange *Exchange)
}

// compileRoute takes Route and returns runtime representation of route.
func compileRoute(c compilerConfig, routeDefinition *Route) (*route, error) {
	producer, err := compileRouteStep(c, routeDefinition.Steps...)
	if err != nil {
		return nil, err
	}

	return &route{
		name:     routeDefinition.Name,
		from:     routeDefinition.From,
		producer: producer,
	}, nil
}

func compileRouteStep(c compilerConfig, s ...RouteStep) (Processor, error) {
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
	// SetBodyStep -> setBodyProcessor
	case *SetBodyStep:
		p := newSetBodyProcessor(t.StepName(), compileExpr(t.BodyValue))
		return decorateProcessor(p, c.preProcessorFunc, c.postProcessorFunc), nil

		// SetHeaderStep -> setHeaderProcessor
	case *SetHeaderStep:
		p := newSetHeaderProcessor(t.StepName(), t.HeaderName, compileExpr(t.HeaderValue))
		return decorateProcessor(p, c.preProcessorFunc, c.postProcessorFunc), nil

		// SetPropertyStep -> setPropertyProcessor
	case *SetPropertyStep:
		p := newSetPropertyProcessor(t.StepName(), t.PropertyName, compileExpr(t.PropertyValue))
		return decorateProcessor(p, c.preProcessorFunc, c.postProcessorFunc), nil

		// ToStep -> toProcessor
	case *ToStep:
		p := newToProcessor(t.StepName(), t.URI)
		return decorateProcessor(p, c.preProcessorFunc, c.postProcessorFunc), nil

		// PipelineStep -> pipelineProcessor
	case *PipelineStep:
		pipeline := newPipelineProcessor(t.StepName(), t.StoOnError)
		for _, step := range t.Steps {
			p, err := compileRouteStep(c, step)
			if err != nil {
				return nil, err
			}
			pipeline.addProcessor(p)
		}
		return decorateProcessor(pipeline, c.preProcessorFunc, c.postProcessorFunc), nil

		// ChoiceStep -> choiceProcessor
	case *ChoiceStep:
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

		// TryStep -> tryProcessor
	case *TryStep:
		p := newTryProcessor(t.StepName())

		for _, tryStep := range t.Steps {
			tp, err := compileRouteStep(c, tryStep)
			if err != nil {
				return nil, err
			}
			p.addProcessor(tp)
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

		// FuncStep -> funcProcessor
	case *FuncStep:
		if userFunc, isUserFunc := t.Func.(func(*Exchange)); isUserFunc {
			return decorateProcessor(newFuncProcessor(t.StepName(), userFunc), c.preProcessorFunc, c.postProcessorFunc), nil
		}
		if storedFuncName, isStoredFunc := t.Func.(string); isStoredFunc {
			storedFunc := c.funcRegistry.Func(storedFuncName)
			if storedFunc == nil {
				return nil, fmt.Errorf("func step: %s: function not found in registry: %s, try to register it first RegisterFunc/MustRegisterFunc",
					t.StepName(), storedFuncName)
			}
			return decorateProcessor(newFuncProcessor(t.StepName(), storedFunc), c.preProcessorFunc, c.postProcessorFunc), nil
		}
		return nil, fmt.Errorf("func step: %s: expected function signature 'func(*Exchange)'", t.StepName())

		// SetErrorStep -> setErrorProcessor
	case *SetErrorStep:
		p := newSetErrorProcessor(t.StepName(), t.Error)
		return decorateProcessor(p, c.preProcessorFunc, c.postProcessorFunc), nil

		// SleepStep -> sleepProcessor
	case *SleepStep:
		p := newSleepProcessor(t.StepName(), t.Duration)
		return decorateProcessor(p, c.preProcessorFunc, c.postProcessorFunc), nil

		// MulticastStep -> multicastProcessor
	case *MulticastStep:
		p := newMulticastProcessor(t.StepName(), t.Parallel, t.StopOnError, t.Aggregator)
		for _, output := range t.Outputs {
			outputProcessor, err := compileRouteStep(c, output.Steps...)
			if err != nil {
				return nil, err
			}
			p.addOutput(outputProcessor)
		}
		return decorateProcessor(p, c.preProcessorFunc, c.postProcessorFunc), nil

	case *LogStep:
		p := newLogProcessor(t.StepName(), t.Msg, t.Level, c.logger)
		return decorateProcessor(p, c.preProcessorFunc, c.postProcessorFunc), nil
	}

	return nil, fmt.Errorf("unknown route step: %T", s[0])
}

func compileExpr(expressionDefinition Expression) expression {
	switch expressionDefinition.Language {
	case "simple":
		s, err := newSimpleExpr(expressionDefinition.Expression)
		if err != nil {
			panic(fmt.Errorf("camel: expression: simple: %w", err))
		}
		return s
	case "constant":
		return newConstExpr(expressionDefinition.Value)
	}

	panic(fmt.Errorf("camel: expression: unknown language: %s", expressionDefinition.Language))
}

func compileErrMatcher(matcherDefinition ErrorMatcher) errorMatcher {
	if strings.TrimSpace(matcherDefinition.Target) == "*" || strings.TrimSpace(matcherDefinition.Target) == "" {
		return errorAny()
	}

	switch matcherDefinition.MatchMode {
	case ErrorMatchModeIs:
		return errorIs(matcherDefinition.Target)
	case ErrorMatchModeContains:
		return errorContains(matcherDefinition.Target)
	case ErrorMatchModeEquals:
		return errorEquals(matcherDefinition.Target)
	case ErrorMatchModeRegex:
		return errorMatches(matcherDefinition.Target)
	}

	panic(fmt.Errorf("camel: error matcher: unknown mode: %s", matcherDefinition.MatchMode))
}
