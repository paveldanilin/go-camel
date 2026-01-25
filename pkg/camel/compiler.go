package camel

import (
	"errors"
	"fmt"
	"github.com/paveldanilin/go-camel/internal/expression"
	"github.com/paveldanilin/go-camel/internal/processor"
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
	"github.com/paveldanilin/go-camel/pkg/camel/route"
	"github.com/paveldanilin/go-camel/pkg/camel/route/errmatcher"
	"github.com/paveldanilin/go-camel/pkg/camel/route/expr"
	"github.com/paveldanilin/go-camel/pkg/camel/route/step"
	"github.com/paveldanilin/go-camel/pkg/camel/template"
	"log/slog"
	"reflect"
	"strings"
)

type compilerConfig struct {
	logger             *slog.Logger
	env                Env
	funcRegistry       FuncRegistry
	dataFormatRegistry DataFormatRegistry
	converterRegistry  ConverterRegistry
	endpointRegistry   EndpointRegistry
	preProcessHook     RouteStepHook
	postProcessHook    RouteStepHook
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

func createProcessor(c compilerConfig, routeName string, steps ...step.RouteStep) (Processor, error) {
	if steps == nil || len(steps) == 0 {
		return nil, errors.New("empty steps")
	}

	if len(steps) > 1 {
		// Auto wrap into pipeline
		pipeline := processor.NewPipeline(false)
		for _, s := range steps {
			p, err := createProcessor(c, routeName, s)
			if err != nil {
				return nil, err
			}
			pipeline.AddProcessor(p)
		}
		return NewRouteStep(routeName, "", pipeline, c.preProcessHook, c.postProcessHook), nil
	}

	var proc Processor
	var err error

	switch t := steps[0].(type) {
	case *step.SetBody:
		proc, err = createSetBodyProcessor(t)

	case *step.SetHeader:
		proc, err = createSetHeaderProcessor(t)

	case *step.SetProperty:
		proc, err = createSetPropertyProcessor(t)

	case *step.To:
		proc, err = createToProcessor(c, t)

	case *step.Pipeline:
		proc, err = createPipelineProcessor(c, routeName, t)

	case *step.Choice:
		proc, err = createChoiceProcessor(c, routeName, t)

	case *step.Try:
		proc, err = createTryProcessor(c, routeName, t)

	case *step.Fn:
		proc, err = createFnProcessor(c, t)

	case *step.SetError:
		proc = processor.NewSetError(t.Error)

	case *step.Delay:
		proc = processor.NewDelay(t.Duration)

	case *step.Multicast:
		proc, err = createMulticastProcessor(c, routeName, t)

	case *step.Log:
		proc, err = processor.NewLog(t.Msg, t.Level, c.logger)

	case *step.RemoveHeader:
		proc = processor.NewRemoveHeader(t.HeaderNames...)

	case *step.RemoveProperty:
		proc = processor.NewRemoveProperty(t.PropertyNames...)

	case *step.Marshal:
		proc = processor.NewMarshal(c.dataFormatRegistry.DataFormat(t.Format))

	case *step.Unmarshal:
		proc = processor.NewUnmarshal(t.TargetType, c.dataFormatRegistry.DataFormat(t.Format))

	case *step.ConvertBody:
		proc, err = createConvertBodyProcessor(c, t)

	case *step.ConvertHeader:
		proc, err = createConvertHeaderProcessor(c, t)

	case *step.ConvertProperty:
		proc, err = createConvertPropertyProcessor(c, t)

	case *step.WhileLoop:
		proc, err = createWhileLoopProcessor(c, routeName, t)

	case *step.ConvertVariable:
		proc, err = createConvertVariableProcessor(c, t)

	case *step.RemoveVariable:
		proc = processor.NewRemoveVariable(t.VariableNames...)

	case *step.SetVariable:
		proc, err = createSetVariableProcessor(t)

	default:
		return nil, fmt.Errorf("unknown route step: %T", steps[0])
	}

	if err != nil {
		return nil, err
	}

	return NewRouteStep(routeName, steps[0].StepName(), proc, c.preProcessHook, c.postProcessHook), nil
}

func createConvertVariableProcessor(c compilerConfig, convertVariable *step.ConvertVariable) (*processor.ConvertVariable, error) {
	var targetType reflect.Type
	if convertVariable.TargetType != nil {
		if reflectedType, isReflectType := convertVariable.TargetType.(reflect.Type); isReflectType {
			targetType = reflectedType
		} else {
			targetType = reflect.TypeOf(convertVariable.TargetType)
		}
	} else if convertVariable.NamedType != "" {
		if namedType, existsNamedType := c.converterRegistry.Type(convertVariable.NamedType); existsNamedType {
			targetType = namedType
		} else {
			return nil, fmt.Errorf("compiler: unknown target type: %s", convertVariable.NamedType)
		}
	} else {
		return nil, errors.New("compiler: no target type")
	}

	return processor.NewConvertVariable(convertVariable.VariableName, targetType, convertVariable.Params, c.converterRegistry), nil
}

func createMulticastProcessor(c compilerConfig, routeName string, multicast *step.Multicast) (*processor.Multicast, error) {
	multicastProcessor := processor.NewMulticast(multicast.Parallel, multicast.StopOnError, multicast.Aggregator)
	for _, output := range multicast.Outputs {
		outputProcessor, err := createProcessor(c, routeName, output.Steps...)
		if err != nil {
			return nil, err
		}
		multicastProcessor.AddProcessor(outputProcessor)
	}
	return multicastProcessor, nil
}

func createWhileLoopProcessor(c compilerConfig, routeName string, whileLoop *step.WhileLoop) (*processor.WhileLoop, error) {
	prdExpr, err := createExpression(whileLoop.Predicate)
	if err != nil {
		return nil, err
	}

	whilePredicate := expression.NewPredicateFromExpression(prdExpr)

	whileLoopProcessor, err := processor.NewWhileLoop(whilePredicate, whileLoop.CopyExchange)
	if err != nil {
		return nil, err
	}

	for _, s := range whileLoop.Steps {
		p, err := createProcessor(c, routeName, s)
		if err != nil {
			return nil, err
		}
		whileLoopProcessor.AddProcessor(p)
	}

	return whileLoopProcessor, nil
}

func createSetBodyProcessor(setBody *step.SetBody) (*processor.SetBody, error) {
	bodyExpr, err := createExpression(setBody.BodyValue)
	if err != nil {
		return nil, err
	}
	return processor.NewSetBody(bodyExpr), nil
}

func createSetHeaderProcessor(setHeader *step.SetHeader) (*processor.SetHeader, error) {
	headerExpr, err := createExpression(setHeader.HeaderValue)
	if err != nil {
		return nil, err
	}
	return processor.NewSetHeader(setHeader.HeaderName, headerExpr), nil
}

func createSetPropertyProcessor(setProperty *step.SetProperty) (*processor.SetProperty, error) {
	propertyExpr, err := createExpression(setProperty.PropertyValue)
	if err != nil {
		return nil, err
	}
	return processor.NewSetProperty(setProperty.PropertyName, propertyExpr), nil
}

func createSetVariableProcessor(setVariable *step.SetVariable) (*processor.SetVariable, error) {
	valueExpr, err := createExpression(setVariable.VariableValue)
	if err != nil {
		return nil, err
	}
	return processor.NewSetVariable(setVariable.VariableName, valueExpr), nil
}

func createToProcessor(c compilerConfig, to *step.To) (Processor, error) {
	uriVars, err := template.Vars(to.URI)
	if err != nil {
		return nil, err
	}

	// static to
	if len(uriVars) == 0 {
		//endpoint := c.endpointRegistry.Endpoint(to.URI)
		//if endpoint == nil {
		//	return nil, fmt.Errorf("failed to create 'to' processor: not found endpoint for URI '%s'", to.URI)
		//}
		//producer, producerErr := endpoint.CreateProducer()
		//if producerErr != nil {
		//	return nil, fmt.Errorf("failed to create 'to' processor: %w", producerErr)
		//}
		producer, err := c.endpointRegistry.CreateEndpointProducer(to.URI)
		if err != nil {
			return nil, err
		}
		return processor.NewTo(producer), nil
	}

	// dynamic to
	uriTpl, uriErr := template.Parse(to.URI)
	if uriErr != nil {
		return nil, fmt.Errorf("failed to create 'to' processor: %w", uriErr)
	}

	return processor.NewToD(uriTpl, uriVars, c.endpointRegistry, c.env), nil
}

func createPipelineProcessor(c compilerConfig, routeName string, pipeline *step.Pipeline) (*processor.Pipeline, error) {
	pipelineProcessor := processor.NewPipeline(pipeline.StoOnError)
	for _, innerStep := range pipeline.Steps {
		p, err := createProcessor(c, routeName, innerStep)
		if err != nil {
			return nil, err
		}
		pipelineProcessor.AddProcessor(p)
	}
	return pipelineProcessor, nil
}

func createChoiceProcessor(c compilerConfig, routeName string, choice *step.Choice) (*processor.Choice, error) {
	choiceProcessor := processor.NewChoice()

	for _, when := range choice.WhenCases {
		thenProcessor, err := createProcessor(c, routeName, when.Steps...)
		if err != nil {
			return nil, fmt.Errorf("failed to create 'then' processor: %w", err)
		}

		prdExpr, err := createExpression(when.Predicate)
		if err != nil {
			return nil, fmt.Errorf("failed to create predicate: %w", err)
		}

		whenPredicate := expression.NewPredicateFromExpression(prdExpr)

		err = choiceProcessor.AddWhen(whenPredicate, thenProcessor)
		if err != nil {
			return nil, fmt.Errorf("failed to add 'when' condition: %w", err)
		}
	}

	if len(choice.Otherwise) > 0 {
		otherwise, err := createProcessor(c, routeName, choice.Otherwise...)
		if err != nil {
			return nil, err
		}
		choiceProcessor.SetOtherwise(otherwise)
	}

	return choiceProcessor, nil
}

func createTryProcessor(c compilerConfig, routeName string, try *step.Try) (*processor.Try, error) {
	tryProcessor := processor.NewTry()

	for _, tryStep := range try.Steps {
		tp, err := createProcessor(c, routeName, tryStep)
		if err != nil {
			return nil, err
		}
		tryProcessor.AddProcessor(tp)
	}

	for _, catch := range try.WhenCatches {
		catchProcessor, err := createProcessor(c, routeName, catch.Steps...)
		if err != nil {
			return nil, err
		}
		tryProcessor.AddCatch(createErrMatcher(catch.ErrorMatcher), catchProcessor)
	}

	for _, finallyStep := range try.FinallySteps {
		finallyProcessor, err := createProcessor(c, routeName, finallyStep)
		if err != nil {
			return nil, err
		}
		tryProcessor.AddFinally(finallyProcessor)
	}

	return tryProcessor, nil
}

func createFnProcessor(c compilerConfig, fn *step.Fn) (*processor.Fn, error) {
	if inlineFunc, isInlineFunc := fn.Func.(func(*exchange.Exchange)); isInlineFunc {
		return processor.NewFn(inlineFunc), nil
	}

	if storedFuncName, isStoredFunc := fn.Func.(string); isStoredFunc {
		storedFunc := c.funcRegistry.Func(storedFuncName)
		if storedFunc == nil {
			return nil, fmt.Errorf("function not found in registry: '%s' at step '%s'",
				storedFuncName, fn.StepName())
		}
		return processor.NewFn(storedFunc), nil
	}

	return nil, fmt.Errorf("unresolved function %v at step '%s'", fn.Func, fn.StepName())
}

func createConvertBodyProcessor(c compilerConfig, convertBody *step.ConvertBody) (*processor.ConvertBody, error) {
	var targetType reflect.Type
	if convertBody.TargetType != nil {
		if reflectedType, isReflectType := convertBody.TargetType.(reflect.Type); isReflectType {
			targetType = reflectedType
		} else {
			targetType = reflect.TypeOf(convertBody.TargetType)
		}
	} else if convertBody.NamedType != "" {
		if namedType, existsNamedType := c.converterRegistry.Type(convertBody.NamedType); existsNamedType {
			targetType = namedType
		} else {
			return nil, fmt.Errorf("compiler: unknown target type: %steps", convertBody.NamedType)
		}
	} else {
		return nil, errors.New("compiler: no target type")
	}

	return processor.NewConvertBody(targetType, convertBody.Params, c.converterRegistry), nil
}

func createConvertHeaderProcessor(c compilerConfig, convertHeader *step.ConvertHeader) (*processor.ConvertHeader, error) {
	var targetType reflect.Type
	if convertHeader.TargetType != nil {
		if reflectedType, isReflectType := convertHeader.TargetType.(reflect.Type); isReflectType {
			targetType = reflectedType
		} else {
			targetType = reflect.TypeOf(convertHeader.TargetType)
		}
	} else if convertHeader.NamedType != "" {
		if namedType, existsNamedType := c.converterRegistry.Type(convertHeader.NamedType); existsNamedType {
			targetType = namedType
		} else {
			return nil, fmt.Errorf("compiler: unknown target type: %steps", convertHeader.NamedType)
		}
	} else {
		return nil, errors.New("compiler: no target type")
	}

	return processor.NewConvertHeader(convertHeader.HeaderName, targetType, convertHeader.Params, c.converterRegistry), nil
}

func createConvertPropertyProcessor(c compilerConfig, convertProperty *step.ConvertProperty) (*processor.ConvertProperty, error) {
	var targetType reflect.Type
	if convertProperty.TargetType != nil {
		if reflectedType, isReflectType := convertProperty.TargetType.(reflect.Type); isReflectType {
			targetType = reflectedType
		} else {
			targetType = reflect.TypeOf(convertProperty.TargetType)
		}
	} else if convertProperty.NamedType != "" {
		if namedType, existsNamedType := c.converterRegistry.Type(convertProperty.NamedType); existsNamedType {
			targetType = namedType
		} else {
			return nil, fmt.Errorf("compiler: unknown target type: %steps", convertProperty.NamedType)
		}
	} else {
		return nil, errors.New("compiler: no target type")
	}

	return processor.NewConvertProperty(convertProperty.PropertyName, targetType, convertProperty.Params, c.converterRegistry), nil
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
		return expression.NewConstant(def.Expression), nil
	case expr.FuncKind:
		if funcExpr, isFuncExpr := def.Expression.(func(e *exchange.Exchange) (any, error)); isFuncExpr {
			return expression.NewFunc(funcExpr), nil
		}
		return nil, fmt.Errorf("failed to create func expression: expected type 'func(e *exchange.Exchange) (any, error)', but got %T", def.Expression)
	}

	return nil, fmt.Errorf("unknown expression kind: %s", def.Kind)
}

func createErrMatcher(matcherDefinition errmatcher.Definition) processor.ErrorMatcher {
	if strings.TrimSpace(matcherDefinition.Target) == "*" || strings.TrimSpace(matcherDefinition.Target) == "" {
		return processor.AnyError()
	}

	switch matcherDefinition.MatchMode {
	case errmatcher.ModeIs:
		return processor.ErrorIs(matcherDefinition.Target)
	case errmatcher.ModeContains:
		return processor.ErrorContains(matcherDefinition.Target)
	case errmatcher.ModeEquals:
		return processor.ErrorEquals(matcherDefinition.Target)
	case errmatcher.ModeRegex:
		return processor.ErrorMatches(matcherDefinition.Target)
	}

	panic(fmt.Errorf("camel: error matcher: unknown mode: %s", matcherDefinition.MatchMode))
}
