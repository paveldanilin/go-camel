package camel

import (
	"fmt"
)

type RouteStep interface {
	StepName() string
}

type Route struct {
	Name  string
	From  string
	Steps []RouteStep
}

// RouteBuilder represents a Route builder.
type RouteBuilder struct {
	route *Route
	stack []*[]RouteStep // step stack
	err   error          // for error tracking
}

func NewRoute(name, from string) *RouteBuilder {
	if name == "" {
		panic(fmt.Errorf("camel: 'name' must be not empty string"))
	}
	if from == "" {
		panic(fmt.Errorf("camel: 'from' must be not empty string"))
	}
	r := &Route{
		Name:  name,
		From:  from,
		Steps: []RouteStep{},
	}
	return &RouteBuilder{
		route: r,
		stack: []*[]RouteStep{&r.Steps},
	}
}

func (b *RouteBuilder) Build() (*Route, error) {
	if b.err != nil {
		return nil, b.err
	}

	// Check size != 1 then not all steps where closed properly.
	if len(b.stack) != 1 {
		return nil, fmt.Errorf("camel: route builder: not all nested steps where closed properly (PipelineStep, ChoiceStep,... etc)")
	}
	return b.route, nil
}

// currentSteps returns the top of the stack.
func (b *RouteBuilder) currentSteps() *[]RouteStep {
	return b.stack[len(b.stack)-1]
}

// addStep adds step to the current context.
func (b *RouteBuilder) addStep(step RouteStep) {
	if b.err != nil {
		return
	}
	steps := b.currentSteps()
	*steps = append(*steps, step)
}

// pushStack pushes builder to the new context.
func (b *RouteBuilder) pushStack(steps *[]RouteStep) {
	b.stack = append(b.stack, steps)
}

// popStack pops builder.
func (b *RouteBuilder) popStack() {
	if len(b.stack) <= 1 {
		b.err = fmt.Errorf("camel: route builder: you are trying to close more block that were opened")
		return
	}
	b.stack = b.stack[:len(b.stack)-1]
}

// ---------------------------------------------------------------------------------------------------------------------
// WalkRoute
// ---------------------------------------------------------------------------------------------------------------------

type StepVisitorFunc func(step RouteStep, depth int) error

func WalkRoute(route *Route, visitor StepVisitorFunc) error {
	if route == nil {
		return nil
	}
	return walkSteps(route.Steps, 0, visitor)
}

func walkSteps(steps []RouteStep, depth int, visitor StepVisitorFunc) error {
	for _, step := range steps {

		if err := visitor(step, depth); err != nil {
			// Stop in case of error
			return err
		}

		switch s := step.(type) {
		case *PipelineStep:
			if err := walkSteps(s.Steps, depth+1, visitor); err != nil {
				return err
			}
		case *ChoiceStep:
			for _, whenCase := range s.WhenCases {
				if err := visitor(&whenCase, depth+1); err != nil {
					return err
				}
				if err := walkSteps(whenCase.Steps, depth+2, visitor); err != nil {
					return err
				}
			}
			if s.Otherwise != nil {
				if err := walkSteps(s.Otherwise, depth+2, visitor); err != nil {
					return err
				}
			}
		case *TryStep:
			// try steps
			if err := walkSteps(s.Steps, depth+1, visitor); err != nil {
				return err
			}
			// catch steps
			for _, catchCase := range s.WhenCatches {
				if err := walkSteps(catchCase.Steps, depth+1, visitor); err != nil {
					return err
				}
			}
			// finally steps
			if s.FinallySteps != nil {
				if err := walkSteps(s.FinallySteps, depth+1, visitor); err != nil {
					return err
				}
			}
		case *LoopWhileStep:
			if err := walkSteps(s.Steps, depth+1, visitor); err != nil {
				return err
			}
		case *MulticastStep:
			for _, subRoute := range s.Outputs {
				if err := visitor(subRoute, depth+1); err != nil {
					return err
				}
				if err := walkSteps(subRoute.Steps, depth+3, visitor); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

type Expression struct {
	Language   string // "simple", "constant",...
	Expression string // Language != "constant"
	Value      any
}

func Simple(expression string) Expression {
	return Expression{
		Language:   "simple",
		Expression: expression,
	}
}

func Constant(value any) Expression {
	return Expression{
		Language: "constant",
		Value:    value,
	}
}
