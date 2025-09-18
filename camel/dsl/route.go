package dsl

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

type PipelineStep struct {
	Name  string
	Steps []RouteStep
}

func (s *PipelineStep) StepName() string { return s.Name }

type SetHeaderStep struct {
	Name        string
	HeaderName  string
	HeaderValue Expression
}

func (s *SetHeaderStep) StepName() string { return s.Name }

type SetBodyStep struct {
	Name      string
	BodyValue Expression
}

func (s *SetBodyStep) StepName() string { return s.Name }

type SetPropertyStep struct {
	Name          string
	PropertyName  string
	PropertyValue Expression
}

func (s *SetPropertyStep) StepName() string { return s.Name }

type ChoiceWhen struct {
	Predicate Expression
	Steps     []RouteStep
}

func (s *ChoiceWhen) StepName() string { return "" }

type ChoiceStep struct {
	Name      string
	WhenCases []ChoiceWhen
	Otherwise []RouteStep
}

func (s *ChoiceStep) StepName() string { return s.Name }

type ToStep struct {
	Name string
	URI  string
}

func (s *ToStep) StepName() string { return s.Name }

type LoopWhileStep struct {
	Name         string
	Predicate    Expression
	CopyExchange bool
	Steps        []RouteStep
}

func (s *LoopWhileStep) StepName() string { return s.Name }

// Do Try

type CatchWhen struct {
	Exception any // (string, func)
	Steps     []RouteStep
}

type DoTryStep struct {
	Name         string
	Steps        []RouteStep
	WhenCatches  []CatchWhen
	FinallySteps []RouteStep
}

func (s *DoTryStep) StepName() string { return s.Name }

type TryStepBuilder struct {
	builder *RouteBuilder
	tryStep *DoTryStep
}

// Catch adds 'catch' to the current DoTryStepKind.
func (tb *TryStepBuilder) Catch(exception any, fn func(b *RouteBuilder)) *TryStepBuilder {
	if tb.builder.err != nil {
		return tb
	}

	catchClause := CatchWhen{Exception: exception}

	tb.builder.pushStack(&catchClause.Steps)
	fn(tb.builder)
	tb.builder.popStack()

	tb.tryStep.WhenCatches = append(tb.tryStep.WhenCatches, catchClause)

	return tb // Catch chain
}

// Finally adds 'finally' and returns to the main builder.
func (tb *TryStepBuilder) Finally(fn func(b *RouteBuilder)) *RouteBuilder {
	if tb.builder.err != nil {
		return tb.builder
	}
	if tb.tryStep.FinallySteps != nil {
		tb.builder.err = fmt.Errorf("block DoTryStepKind '%s' already has block Finally", tb.tryStep.Name)
		return tb.builder
	}

	tb.builder.pushStack(&tb.tryStep.FinallySteps)
	fn(tb.builder)
	tb.builder.popStack()

	return tb.builder // main builder
}

func (tb *TryStepBuilder) EndTry() *RouteBuilder {
	return tb.builder
}

// ---------------------------------------------------------------------------------------------------------------------
// RouteBuilder
// ---------------------------------------------------------------------------------------------------------------------

// RouteBuilder represents a Route builder.
type RouteBuilder struct {
	route      *Route
	stack      []*[]RouteStep // step stack
	lastChoice *ChoiceStep    // link to the last Choice for adding When/Otherwise.
	err        error          // for error tracking
}

func NewRouteBuilder(name, from string) *RouteBuilder {
	route := &Route{
		Name:  name,
		From:  from,
		Steps: []RouteStep{},
	}
	return &RouteBuilder{
		route: route,
		stack: []*[]RouteStep{&route.Steps},
	}
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
		b.err = fmt.Errorf("ошибка конфигурации: попытка закрыть больше блоков, чем было открыто")
		return
	}
	b.stack = b.stack[:len(b.stack)-1]
}

// SetHeader adds steps to set message header.
func (b *RouteBuilder) SetHeader(stepName, headerName string, headerValue Expression) *RouteBuilder {
	if b.err != nil {
		return b
	}
	b.addStep(&SetHeaderStep{
		Name:        stepName,
		HeaderName:  headerName,
		HeaderValue: headerValue,
	})
	return b
}

// SetBody adds step to set the message body.
func (b *RouteBuilder) SetBody(stepName string, bodyValue Expression) *RouteBuilder {
	if b.err != nil {
		return b
	}
	b.addStep(&SetBodyStep{
		Name:      stepName,
		BodyValue: bodyValue,
	})
	return b
}

// SetProperty adds step to set a exchange property.
func (b *RouteBuilder) SetProperty(stepName, propertyName string, propertyValue Expression) *RouteBuilder {
	if b.err != nil {
		return b
	}
	b.addStep(&SetPropertyStep{
		Name:          stepName,
		PropertyName:  propertyName,
		PropertyValue: propertyValue,
	})
	return b
}

// Pipeline adds new pipeline.
// Function configure will be called to configure PipelineStep.
func (b *RouteBuilder) Pipeline(stepName string, configure func(b *RouteBuilder)) *RouteBuilder {
	if b.err != nil {
		return b
	}
	pipeline := &PipelineStep{Name: stepName, Steps: []RouteStep{}}
	b.addStep(pipeline)

	// push pipeline
	b.pushStack(&pipeline.Steps)
	configure(b) // configure PipelineStep
	b.popStack() // pop

	return b
}

// Choice adds choice step.
// Function configure will be called to configure ChoiceStep.
func (b *RouteBuilder) Choice(stepName string, configure func(b *RouteBuilder)) *RouteBuilder {
	if b.err != nil {
		return b
	}
	choice := &ChoiceStep{Name: stepName}
	b.addStep(choice)

	// store link to the current choice
	prevChoice := b.lastChoice
	b.lastChoice = choice

	configure(b) // configure WhenStepKind/Otherwise

	// restore link to the external choice (if any)
	b.lastChoice = prevChoice
	return b
}

// When добавляет условную ветку в последний созданный ChoiceStepKind.
func (b *RouteBuilder) When(predicate Expression, configure func(b *RouteBuilder)) *RouteBuilder {
	if b.err != nil {
		return b
	}
	if b.lastChoice == nil {
		b.err = fmt.Errorf("bad method call: WhenStepKind() must be called only inside ChoiceStepKind()")
		return b
	}

	whenCase := ChoiceWhen{Predicate: predicate, Steps: []RouteStep{}}
	b.lastChoice.WhenCases = append(b.lastChoice.WhenCases, whenCase)

	stepsPtr := &b.lastChoice.WhenCases[len(b.lastChoice.WhenCases)-1].Steps

	b.pushStack(stepsPtr)
	configure(b) // configure WhenStepKind
	b.popStack()

	return b
}

func (b *RouteBuilder) Otherwise(configure func(b *RouteBuilder)) *RouteBuilder {
	if b.err != nil {
		return b
	}
	if b.lastChoice == nil {
		b.err = fmt.Errorf("bad method call: Otherwise() must be called only inside ChoiceStepKind()")
		return b
	}
	if b.lastChoice.Otherwise != nil {
		b.err = fmt.Errorf("could not redefine Otherwise since it already created for ChoiceStepKind")
		return b
	}

	b.lastChoice.Otherwise = []RouteStep{}
	b.pushStack(&b.lastChoice.Otherwise)
	configure(b) // configure Otherwise
	b.popStack()

	return b
}

// DoTry adds try-catch-finally step.
func (b *RouteBuilder) DoTry(stepName string, configure func(b *RouteBuilder)) *TryStepBuilder {
	if b.err != nil {
		return &TryStepBuilder{builder: b}
	}

	tryStep := &DoTryStep{Name: stepName}
	b.addStep(tryStep)

	// configure try
	b.pushStack(&tryStep.Steps)
	configure(b)
	b.popStack()

	return &TryStepBuilder{builder: b, tryStep: tryStep}
}

func (b *RouteBuilder) To(stepName, uri string) *RouteBuilder {
	if b.err != nil {
		return b
	}
	b.addStep(&ToStep{
		Name: stepName,
		URI:  uri,
	})

	return b
}

func (b *RouteBuilder) LoopWhile(stepName string, predicate Expression, copyExchange bool, configure func(b *RouteBuilder)) *RouteBuilder {
	if b.err != nil {
		return b
	}

	step := &LoopWhileStep{
		Name:         stepName,
		Predicate:    predicate,
		CopyExchange: copyExchange,
	}
	b.addStep(step)

	b.pushStack(&step.Steps)
	configure(b)
	b.popStack()

	return b
}

func (b *RouteBuilder) Build() (*Route, error) {
	if b.err != nil {
		return nil, b.err
	}

	// Check size != 1 then not all steps where closed properly.
	if len(b.stack) != 1 {
		return nil, fmt.Errorf("bad configuration: not all nested steps where closed properly (PipelineStepKind, ChoiceStepKind,... etc)")
	}
	return b.route, nil
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
		case *DoTryStep:
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
		}
	}
	return nil
}
