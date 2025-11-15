package camel

import (
	"errors"
	"fmt"
	"github.com/paveldanilin/go-camel/internal/eip/choice"
	"github.com/paveldanilin/go-camel/internal/eip/convertbody"
	"github.com/paveldanilin/go-camel/internal/eip/convertheader"
	"github.com/paveldanilin/go-camel/internal/eip/convertproperty"
	"github.com/paveldanilin/go-camel/internal/eip/delay"
	"github.com/paveldanilin/go-camel/internal/eip/fn"
	"github.com/paveldanilin/go-camel/internal/eip/log"
	"github.com/paveldanilin/go-camel/internal/eip/marshal"
	"github.com/paveldanilin/go-camel/internal/eip/multicast"
	"github.com/paveldanilin/go-camel/internal/eip/pipeline"
	"github.com/paveldanilin/go-camel/internal/eip/removeheader"
	"github.com/paveldanilin/go-camel/internal/eip/removeproperty"
	"github.com/paveldanilin/go-camel/internal/eip/setbody"
	"github.com/paveldanilin/go-camel/internal/eip/seterror"
	"github.com/paveldanilin/go-camel/internal/eip/setheader"
	"github.com/paveldanilin/go-camel/internal/eip/setproperty"
	"github.com/paveldanilin/go-camel/internal/eip/to"
	"github.com/paveldanilin/go-camel/internal/eip/try"
	"github.com/paveldanilin/go-camel/internal/eip/unmarshal"
	"github.com/paveldanilin/go-camel/internal/expression"
	"github.com/paveldanilin/go-camel/internal/processor"
	"github.com/paveldanilin/go-camel/pkg/camel/api"
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
	"github.com/paveldanilin/go-camel/pkg/camel/route"
	"github.com/paveldanilin/go-camel/pkg/camel/route/errmatcher"
	"github.com/paveldanilin/go-camel/pkg/camel/route/expr"
	"github.com/paveldanilin/go-camel/pkg/camel/route/step"
	"github.com/paveldanilin/go-camel/pkg/camel/template"
	"reflect"
	"strings"
)

type compilerConfig struct {
	logger             api.Logger
	env                api.Env
	funcRegistry       FuncRegistry
	dataFormatRegistry DataFormatRegistry
	converterRegistry  ConverterRegistry
	endpointRegistry   EndpointRegistry
	preProcessor       func(e *exchange.Exchange)
	postProcessor      func(e *exchange.Exchange)
}

// compileRoute converts a set of route.Definition to a set of api.Processor.
func compileRoute(c compilerConfig, routeDefinition *route.Definition) (*rtRoute, error) {
	producer, err := createProcessor(c, routeDefinition.Name, routeDefinition.Steps...)
	if err != nil {
		return nil, err
	}

	return &rtRoute{
		name:     routeDefinition.Name,
		from:     routeDefinition.From,
		producer: producer,
	}, nil
}

func createProcessor(c compilerConfig, routeName string, s ...api.RouteStep) (api.Processor, error) {
	if s == nil || len(s) == 0 {
		return nil, errors.New("empty steps")
	}

	if len(s) > 1 {
		pipe := pipeline.NewProcessor(routeName, "", false)
		for _, innerStep := range s {
			p, err := createProcessor(c, routeName, innerStep)
			if err != nil {
				return nil, err
			}
			pipe.AddProcessor(p)
		}
		return processor.NewProcessor(pipe, c.preProcessor, c.postProcessor), nil
	}

	switch t := s[0].(type) {
	case *step.SetBody:
		bodyExpr, err := createExpression(t.BodyValue)
		if err != nil {
			return nil, err
		}
		p := setbody.NewProcessor(routeName, t.StepName(), bodyExpr)
		return processor.NewProcessor(p, c.preProcessor, c.postProcessor), nil

	case *step.SetHeader:
		headerExpr, err := createExpression(t.HeaderValue)
		if err != nil {
			return nil, err
		}
		p := setheader.NewProcessor(routeName, t.StepName(), t.HeaderName, headerExpr)
		return processor.NewProcessor(p, c.preProcessor, c.postProcessor), nil

	case *step.SetProperty:
		propertyExpr, err := createExpression(t.PropertyValue)
		if err != nil {
			return nil, err
		}
		p := setproperty.NewProcessor(routeName, t.StepName(), t.PropertyName, propertyExpr)
		return processor.NewProcessor(p, c.preProcessor, c.postProcessor), nil

	case *step.To:
		uriVars, err := template.Vars(t.URI)
		if err != nil {
			return nil, err
		}

		var p api.Processor
		if len(uriVars) == 0 {
			endpoint := c.endpointRegistry.Endpoint(t.URI)
			if endpoint == nil {
				return nil, fmt.Errorf("failed to create 'to' processor: not found endpoint for URI '%s'", t.URI)
			}
			producer, producerErr := endpoint.CreateProducer()
			if producerErr != nil {
				return nil, fmt.Errorf("failed to create 'to' processor: %w", producerErr)
			}
			p = to.NewStaticProcessor(routeName, t.StepName(), producer)
		} else {
			uriTpl, uriErr := template.Parse(t.URI)
			if uriErr != nil {
				return nil, fmt.Errorf("failed to create 'to' processor: %w", uriErr)
			}

			p = to.NewDynamicProcessor(routeName, t.StepName(), uriTpl, uriVars, c.endpointRegistry, c.env)
		}
		return processor.NewProcessor(p, c.preProcessor, c.postProcessor), nil

	case *step.Pipeline:
		pipe := pipeline.NewProcessor(routeName, t.StepName(), t.StoOnError)
		for _, innerStep := range t.Steps {
			p, err := createProcessor(c, routeName, innerStep)
			if err != nil {
				return nil, err
			}
			pipe.AddProcessor(p)
		}
		return processor.NewProcessor(pipe, c.preProcessor, c.postProcessor), nil

	case *step.Choice:
		p := choice.NewProcessor(routeName, t.StepName())
		for _, when := range t.WhenCases {
			whenBody, err := createProcessor(c, routeName, when.Steps...)
			if err != nil {
				return nil, err
			}

			prdExpr, err := createExpression(when.Predicate)
			if err != nil {
				return nil, err
			}
			p.AddWhen(prdExpr, whenBody)
		}
		if len(t.Otherwise) > 0 {
			otherwise, err := createProcessor(c, routeName, t.Otherwise...)
			if err != nil {
				return nil, err
			}
			p.SetOtherwise(otherwise)
		}
		return processor.NewProcessor(p, c.preProcessor, c.postProcessor), nil

	case *step.Try:
		p := try.NewProcessor(routeName, t.StepName())

		for _, tryStep := range t.Steps {
			tp, err := createProcessor(c, routeName, tryStep)
			if err != nil {
				return nil, err
			}
			p.AddProcessor(tp)
		}

		for _, catch := range t.WhenCatches {
			catchProcessor, err := createProcessor(c, routeName, catch.Steps...)
			if err != nil {
				return nil, err
			}
			p.AddCatch(createErrMatcher(catch.ErrorMatcher), catchProcessor)
		}

		for _, finallyStep := range t.FinallySteps {
			finallyProcessor, err := createProcessor(c, routeName, finallyStep)
			if err != nil {
				return nil, err
			}
			p.AddFinally(finallyProcessor)
		}
		return processor.NewProcessor(p, c.preProcessor, c.postProcessor), nil

	case *step.Fn:
		if inlineFunc, isInlineFunc := t.Func.(func(*exchange.Exchange)); isInlineFunc {
			return processor.NewProcessor(fn.NewProcessor(routeName, t.StepName(), inlineFunc), c.preProcessor, c.postProcessor), nil
		}
		if storedFuncName, isStoredFunc := t.Func.(string); isStoredFunc {
			storedFunc := c.funcRegistry.Func(storedFuncName)
			if storedFunc == nil {
				return nil, fmt.Errorf("fn step: %s: function not found in registry: %s, try to register it first RegisterFunc/MustRegisterFunc",
					t.StepName(), storedFuncName)
			}
			return processor.NewProcessor(fn.NewProcessor(routeName, t.StepName(), storedFunc), c.preProcessor, c.postProcessor), nil
		}
		return nil, fmt.Errorf("fn step: %s: expected function signature 'fn(*Exchange)'", t.StepName())

	case *step.SetError:
		p := seterror.NewProcessor(routeName, t.StepName(), t.Error)
		return processor.NewProcessor(p, c.preProcessor, c.postProcessor), nil

	case *step.Delay:
		p := delay.NewProcessor(routeName, t.StepName(), t.Duration)
		return processor.NewProcessor(p, c.preProcessor, c.postProcessor), nil

	case *step.Multicast:
		p := multicast.NewProcessor(routeName, t.StepName(), t.Parallel, t.StopOnError, t.Aggregator)
		for _, output := range t.Outputs {
			outputProcessor, err := createProcessor(c, routeName, output.Steps...)
			if err != nil {
				return nil, err
			}
			p.AddOutput(outputProcessor)
		}
		return processor.NewProcessor(p, c.preProcessor, c.postProcessor), nil

	case *step.Log:
		p := log.NewProcessor(routeName, t.StepName(), t.Msg, t.Level, c.logger)
		return processor.NewProcessor(p, c.preProcessor, c.postProcessor), nil

	case *step.RemoveHeader:
		p := removeheader.NewProcessor(routeName, t.StepName(), t.HeaderNames...)
		return processor.NewProcessor(p, c.preProcessor, c.postProcessor), nil

	case *step.RemoveProperty:
		p := removeproperty.NewProcessor(routeName, t.StepName(), t.PropertyNames...)
		return processor.NewProcessor(p, c.preProcessor, c.postProcessor), nil

	case *step.Marshal:
		p := marshal.NewProcessor(routeName, t.StepName(), c.dataFormatRegistry.DataFormat(t.Format))
		return processor.NewProcessor(p, c.preProcessor, c.postProcessor), nil

	case *step.Unmarshal:
		p := unmarshal.NewProcessor(routeName, t.StepName(), t.TargetType, c.dataFormatRegistry.DataFormat(t.Format))
		return processor.NewProcessor(p, c.preProcessor, c.postProcessor), nil

	case *step.ConvertBody:
		var targetType reflect.Type
		if t.TargetType != nil {
			if reflectedType, isReflectType := t.TargetType.(reflect.Type); isReflectType {
				targetType = reflectedType
			} else {
				targetType = reflect.TypeOf(t.TargetType)
			}
		} else if t.NamedType != "" {
			if namedType, existsNamedType := c.converterRegistry.Type(t.NamedType); existsNamedType {
				targetType = namedType
			} else {
				return nil, fmt.Errorf("compiler: unknown target type: %s", t.NamedType)
			}
		} else {
			return nil, errors.New("compiler: no target type")
		}

		p := convertbody.NewProcessor(routeName, t.StepName(), targetType, t.Params, c.converterRegistry)
		return processor.NewProcessor(p, c.preProcessor, c.postProcessor), nil

	case *step.ConvertHeader:
		var targetType reflect.Type
		if t.TargetType != nil {
			if reflectedType, isReflectType := t.TargetType.(reflect.Type); isReflectType {
				targetType = reflectedType
			} else {
				targetType = reflect.TypeOf(t.TargetType)
			}
		} else if t.NamedType != "" {
			if namedType, existsNamedType := c.converterRegistry.Type(t.NamedType); existsNamedType {
				targetType = namedType
			} else {
				return nil, fmt.Errorf("compiler: unknown target type: %s", t.NamedType)
			}
		} else {
			return nil, errors.New("compiler: no target type")
		}

		p := convertheader.NewProcessor(routeName, t.StepName(), t.HeaderName, targetType, t.Params, c.converterRegistry)
		return processor.NewProcessor(p, c.preProcessor, c.postProcessor), nil

	case *step.ConvertProperty:
		var targetType reflect.Type
		if t.TargetType != nil {
			if reflectedType, isReflectType := t.TargetType.(reflect.Type); isReflectType {
				targetType = reflectedType
			} else {
				targetType = reflect.TypeOf(t.TargetType)
			}
		} else if t.NamedType != "" {
			if namedType, existsNamedType := c.converterRegistry.Type(t.NamedType); existsNamedType {
				targetType = namedType
			} else {
				return nil, fmt.Errorf("compiler: unknown target type: %s", t.NamedType)
			}
		} else {
			return nil, errors.New("compiler: no target type")
		}

		p := convertproperty.NewProcessor(routeName, t.StepName(), t.PropertyName, targetType, t.Params, c.converterRegistry)
		return processor.NewProcessor(p, c.preProcessor, c.postProcessor), nil
	}

	return nil, fmt.Errorf("unknown route step: %T", s[0])
}

func createExpression(def expr.Definition) (expression.Expression, error) {
	switch def.Kind {
	case expr.SimpleKind:
		se, err := expression.NewSimple(def.Expression.(string))
		if err != nil {
			return nil, fmt.Errorf("failed to create simple expression: %w", err)
		}
		return se, nil
	case expr.ConstantKind:
		return expression.NewConst(def.Expression), nil
	case expr.FuncKind:
		if funcExpr, isFuncExpr := def.Expression.(func(e *exchange.Exchange) (any, error)); isFuncExpr {
			return expression.NewFunc(funcExpr), nil
		}
		return nil, fmt.Errorf("failed to create func expression: expected type 'func(e *exchange.Exchange) (any, error)', but got %T", def.Expression)
	}

	return nil, fmt.Errorf("unknown expression kind: %s", def.Kind)
}

func createErrMatcher(matcherDefinition errmatcher.Definition) try.ErrorMatcher {
	if strings.TrimSpace(matcherDefinition.Target) == "*" || strings.TrimSpace(matcherDefinition.Target) == "" {
		return try.AnyError()
	}

	switch matcherDefinition.MatchMode {
	case errmatcher.ModeIs:
		return try.ErrorIs(matcherDefinition.Target)
	case errmatcher.ModeContains:
		return try.ErrorContains(matcherDefinition.Target)
	case errmatcher.ModeEquals:
		return try.ErrorEquals(matcherDefinition.Target)
	case errmatcher.ModeRegex:
		return try.ErrorMatches(matcherDefinition.Target)
	}

	panic(fmt.Errorf("camel: error matcher: unknown mode: %s", matcherDefinition.MatchMode))
}
