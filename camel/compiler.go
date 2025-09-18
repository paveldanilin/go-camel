package camel

import (
	"errors"
	"fmt"
	"github.com/paveldanilin/go-camel/camel/dsl"
)

type compilerConfig struct {
	preProcessorFunc  func(exchange *Exchange)
	postProcessorFunc func(exchange *Exchange)
}

// compileRoute takes dsl.Route and returns runtime representation of route -> route.
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
		pipeline := newPipelineProcessor()
		for _, step := range s {
			p, err := compileRouteStep(c, step)
			if err != nil {
				return nil, err
			}
			pipeline.WithProcessor(p)
		}
		return newProcessor(pipeline, c.preProcessorFunc, c.postProcessorFunc), nil
	}

	switch t := s[0].(type) {
	case *dsl.SetBodyStep:
		p := newSetBodyProcessor(compileExpr(t.BodyValue)).WithStepName(t.StepName())
		return newProcessor(p, c.preProcessorFunc, c.postProcessorFunc), nil
	case *dsl.SetHeaderStep:
		p := newSetHeaderProcessor(t.HeaderName, compileExpr(t.HeaderValue)).WithStepName(t.StepName())
		return newProcessor(p, c.preProcessorFunc, c.postProcessorFunc), nil
	case *dsl.ToStep:
		p := newToProcessor(t.URI).WithStepName(t.StepName())
		return newProcessor(p, c.preProcessorFunc, c.postProcessorFunc), nil
	case *dsl.PipelineStep:
		pipeline := newPipelineProcessor().WithStepName(t.StepName())
		for _, step := range t.Steps {
			p, err := compileRouteStep(c, step)
			if err != nil {
				return nil, err
			}
			pipeline.WithProcessor(p)
		}
		return newProcessor(pipeline, c.preProcessorFunc, c.postProcessorFunc), nil
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
