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
	preProcessor       func(exchange *Exchange)
	postProcessor      func(exchange *Exchange)
}

// compileRoute takes Route definition and returns runtime representation of the route.
func compileRoute(c compilerConfig, routeDefinition *Route) (*route, error) {
	producer, err := createProcessor(c, routeDefinition.Name, routeDefinition.Steps...)
	if err != nil {
		return nil, err
	}

	return &route{
		name:     routeDefinition.Name,
		from:     routeDefinition.From,
		producer: producer,
	}, nil
}

func createProcessor(c compilerConfig, routeName string, s ...RouteStep) (Processor, error) {
	if s == nil || len(s) == 0 {
		return nil, errors.New("empty steps")
	}

	if len(s) > 1 {
		pipeline := newPipelineProcessor("", false)
		for _, step := range s {
			p, err := createProcessor(c, routeName, step)
			if err != nil {
				return nil, err
			}
			pipeline.addProcessor(p)
		}
		return decorateProcessor(pipeline, c.preProcessor, c.postProcessor), nil
	}

	switch t := s[0].(type) {
	// SetBodyStep -> setBodyProcessor
	case *SetBodyStep:
		p := newSetBodyProcessor(t.StepName(), createExpression(t.BodyValue))
		return decorateProcessor(p, c.preProcessor, c.postProcessor), nil

		// SetHeaderStep -> setHeaderProcessor
	case *SetHeaderStep:
		p := newSetHeaderProcessor(t.StepName(), t.HeaderName, createExpression(t.HeaderValue))
		return decorateProcessor(p, c.preProcessor, c.postProcessor), nil

		// SetPropertyStep -> setPropertyProcessor
	case *SetPropertyStep:
		p := newSetPropertyProcessor(t.StepName(), t.PropertyName, createExpression(t.PropertyValue))
		return decorateProcessor(p, c.preProcessor, c.postProcessor), nil

		// ToStep -> toProcessor
	case *ToStep:
		p := newToProcessor(t.StepName(), t.URI)
		return decorateProcessor(p, c.preProcessor, c.postProcessor), nil

		// PipelineStep -> pipelineProcessor
	case *PipelineStep:
		pipeline := newPipelineProcessor(t.StepName(), t.StoOnError)
		for _, step := range t.Steps {
			p, err := createProcessor(c, routeName, step)
			if err != nil {
				return nil, err
			}
			pipeline.addProcessor(p)
		}
		return decorateProcessor(pipeline, c.preProcessor, c.postProcessor), nil

		// ChoiceStep -> choiceProcessor
	case *ChoiceStep:
		p := newChoiceProcessor(t.StepName())
		for _, when := range t.WhenCases {
			whenBody, err := createProcessor(c, routeName, when.Steps...)
			if err != nil {
				return nil, err
			}
			p.addWhen(createExpression(when.Predicate), whenBody)
		}
		if len(t.Otherwise) > 0 {
			otherwise, err := createProcessor(c, routeName, t.Otherwise...)
			if err != nil {
				return nil, err
			}
			p.setOtherwise(otherwise)
		}
		return decorateProcessor(p, c.preProcessor, c.postProcessor), nil

		// TryStep -> tryProcessor
	case *TryStep:
		p := newTryProcessor(t.StepName())

		for _, tryStep := range t.Steps {
			tp, err := createProcessor(c, routeName, tryStep)
			if err != nil {
				return nil, err
			}
			p.addProcessor(tp)
		}

		for _, catch := range t.WhenCatches {
			catchProcessor, err := createProcessor(c, routeName, catch.Steps...)
			if err != nil {
				return nil, err
			}
			p.addCatch(createErrMatcher(catch.ErrorMatcher), catchProcessor)
		}

		for _, finallyStep := range t.FinallySteps {
			finallyProcessor, err := createProcessor(c, routeName, finallyStep)
			if err != nil {
				return nil, err
			}
			p.addFinally(finallyProcessor)
		}
		return decorateProcessor(p, c.preProcessor, c.postProcessor), nil

		// FuncStep -> funcProcessor
	case *FuncStep:
		if userFunc, isUserFunc := t.Func.(func(*Exchange)); isUserFunc {
			return decorateProcessor(newFuncProcessor(t.StepName(), userFunc), c.preProcessor, c.postProcessor), nil
		}
		if storedFuncName, isStoredFunc := t.Func.(string); isStoredFunc {
			storedFunc := c.funcRegistry.Func(storedFuncName)
			if storedFunc == nil {
				return nil, fmt.Errorf("func step: %s: function not found in registry: %s, try to register it first RegisterFunc/MustRegisterFunc",
					t.StepName(), storedFuncName)
			}
			return decorateProcessor(newFuncProcessor(t.StepName(), storedFunc), c.preProcessor, c.postProcessor), nil
		}
		return nil, fmt.Errorf("func step: %s: expected function signature 'func(*Exchange)'", t.StepName())

		// SetErrorStep -> setErrorProcessor
	case *SetErrorStep:
		p := newSetErrorProcessor(t.StepName(), t.Error)
		return decorateProcessor(p, c.preProcessor, c.postProcessor), nil

		// SleepStep -> sleepProcessor
	case *SleepStep:
		p := newSleepProcessor(t.StepName(), t.Duration)
		return decorateProcessor(p, c.preProcessor, c.postProcessor), nil

		// MulticastStep -> multicastProcessor
	case *MulticastStep:
		p := newMulticastProcessor(t.StepName(), t.Parallel, t.StopOnError, t.Aggregator)
		for _, output := range t.Outputs {
			outputProcessor, err := createProcessor(c, routeName, output.Steps...)
			if err != nil {
				return nil, err
			}
			p.addOutput(outputProcessor)
		}
		return decorateProcessor(p, c.preProcessor, c.postProcessor), nil

	case *LogStep:
		p := newLogProcessor(routeName, t.StepName(), t.Msg, t.Level, c.logger)
		return decorateProcessor(p, c.preProcessor, c.postProcessor), nil

	case *RemoveHeaderStep:
		p := newRemoveHeaderProcessor(t.StepName(), t.HeaderName)
		return decorateProcessor(p, c.preProcessor, c.postProcessor), nil

	case *RemovePropertyStep:
		p := newRemovePropertyProcessor(t.StepName(), t.PropertyName)
		return decorateProcessor(p, c.preProcessor, c.postProcessor), nil

	case *MarshalStep:
		p := newMarshalProcessor(t.StepName(), c.dataFormatRegistry.DataFormat(t.Format))
		return decorateProcessor(p, c.preProcessor, c.postProcessor), nil

	case *UnmarshalStep:
		p := newUnmarshalProcessor(t.StepName(), t.TargetType, c.dataFormatRegistry.DataFormat(t.Format))
		return decorateProcessor(p, c.preProcessor, c.postProcessor), nil
	}

	return nil, fmt.Errorf("unknown route step: %T", s[0])
}

func createExpression(expressionDefinition Expression) expression {
	switch expressionDefinition.Language {
	case "simple":
		se, err := newSimpleExpression(expressionDefinition.Expression.(string))
		if err != nil {
			panic(fmt.Errorf("camel: expression: simple: %w", err))
		}
		return se
	case "constant":
		return newConstExpression(expressionDefinition.Expression)
	}

	panic(fmt.Errorf("camel: expression: unknown language: %s", expressionDefinition.Language))
}

func createErrMatcher(matcherDefinition ErrorMatcher) errorMatcher {
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
