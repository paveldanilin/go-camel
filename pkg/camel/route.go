package camel

import (
	"fmt"
	"github.com/paveldanilin/go-camel/pkg/camel/api"
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
	"github.com/paveldanilin/go-camel/pkg/camel/expr"
	"github.com/paveldanilin/go-camel/pkg/camel/routestep"
)

type Route struct {
	Name  string
	From  string
	Steps []api.RouteStep
}

// RouteBuilder represents a Route builder.
type RouteBuilder struct {
	route *Route
	stack []*[]api.RouteStep
	err   error
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
		Steps: []api.RouteStep{},
	}
	return &RouteBuilder{
		route: r,
		stack: []*[]api.RouteStep{&r.Steps},
	}
}

func (b *RouteBuilder) Build() (*Route, error) {
	if b.err != nil {
		return nil, b.err
	}

	// Check size != 1 then not all steps where closed properly.
	if len(b.stack) != 1 {
		return nil, fmt.Errorf("camel: step builder: not all nested steps where closed properly (PipelineStep, ChoiceStep,... etc)")
	}
	return b.route, nil
}

// currentSteps returns the top of the stack.
func (b *RouteBuilder) currentSteps() *[]api.RouteStep {
	return b.stack[len(b.stack)-1]
}

// addStep adds step to the current context.
func (b *RouteBuilder) addStep(step api.RouteStep) {
	if b.err != nil {
		return
	}
	steps := b.currentSteps()
	*steps = append(*steps, step)
}

// pushStack pushes builder to the new context.
func (b *RouteBuilder) pushStack(steps *[]api.RouteStep) {
	b.stack = append(b.stack, steps)
}

// popStack pops builder.
func (b *RouteBuilder) popStack() {
	if len(b.stack) <= 1 {
		b.err = fmt.Errorf("camel: step builder: you are trying to close more block that were opened")
		return
	}
	b.stack = b.stack[:len(b.stack)-1]
}

// SetBody adds step to set the message body.
func (b *RouteBuilder) SetBody(stepName string, bodyValue expr.Expression) *RouteBuilder {
	if b.err != nil {
		return b
	}
	b.addStep(&routestep.SetBody{
		Name:      stepName,
		BodyValue: bodyValue,
	})
	return b
}

func (b *RouteBuilder) ConvertBodyTo(stepName string, targetType any, params map[string]any) *RouteBuilder {
	if b.err != nil {
		return b
	}
	b.addStep(&routestep.ConvertBody{
		Name:       stepName,
		TargetType: targetType,
		Params:     params,
	})
	return b
}

func (b *RouteBuilder) ConvertBodyToNamed(stepName, typeName string, params map[string]any) *RouteBuilder {
	if b.err != nil {
		return b
	}
	b.addStep(&routestep.ConvertBody{
		Name:      stepName,
		NamedType: typeName,
		Params:    params,
	})
	return b
}

func (b *RouteBuilder) ConvertHeaderTo(stepName, headerName string, targetType any, params map[string]any) *RouteBuilder {
	if b.err != nil {
		return b
	}
	b.addStep(&routestep.ConvertHeader{
		Name:       stepName,
		TargetType: targetType,
		HeaderName: headerName,
		Params:     params,
	})
	return b
}

func (b *RouteBuilder) ConvertHeaderToNamed(stepName, headerName, typeName string, params map[string]any) *RouteBuilder {
	if b.err != nil {
		return b
	}
	b.addStep(&routestep.ConvertHeader{
		Name:       stepName,
		NamedType:  typeName,
		HeaderName: headerName,
		Params:     params,
	})
	return b
}

func (b *RouteBuilder) ConvertPropertyTo(stepName, propertyName string, targetType any, params map[string]any) *RouteBuilder {
	if b.err != nil {
		return b
	}
	b.addStep(&routestep.ConvertProperty{
		Name:         stepName,
		TargetType:   targetType,
		PropertyName: propertyName,
		Params:       params,
	})
	return b
}

func (b *RouteBuilder) ConvertPropertyToNamed(stepName, propertyName, typeName string, params map[string]any) *RouteBuilder {
	if b.err != nil {
		return b
	}
	b.addStep(&routestep.ConvertProperty{
		Name:         stepName,
		NamedType:    typeName,
		PropertyName: propertyName,
		Params:       params,
	})
	return b
}

// Choice adds choice step to the current step level.
func (b *RouteBuilder) Choice(stepName string) *ChoiceStepBuilder {
	if b.err != nil {
		return &ChoiceStepBuilder{builder: b}
	}

	choiceStep := &routestep.Choice{Name: stepName}
	b.addStep(choiceStep)

	return &ChoiceStepBuilder{builder: b, choiceStep: choiceStep}
}

// SetHeader adds steps to set message header.
func (b *RouteBuilder) SetHeader(stepName, headerName string, headerValue expr.Expression) *RouteBuilder {
	if b.err != nil {
		return b
	}
	b.addStep(&routestep.SetHeader{
		Name:        stepName,
		HeaderName:  headerName,
		HeaderValue: headerValue,
	})
	return b
}

// Try adds try-catch-finally step.
func (b *RouteBuilder) Try(stepName string, configure func(b *RouteBuilder)) *TryStepBuilder {
	if b.err != nil {
		return &TryStepBuilder{builder: b}
	}

	tryStep := &routestep.Try{Name: stepName}
	b.addStep(tryStep)

	b.pushStack(&tryStep.Steps)
	configure(b) // configure try
	b.popStack()

	return &TryStepBuilder{builder: b, tryStep: tryStep}
}

func (b *RouteBuilder) SetError(stepName string, err error) *RouteBuilder {
	if b.err != nil {
		return b
	}
	b.addStep(&routestep.SetError{
		Name:  stepName,
		Error: err,
	})
	return b
}

// SetProperty adds step to set an exchange property.
func (b *RouteBuilder) SetProperty(stepName, propertyName string, propertyValue expr.Expression) *RouteBuilder {
	if b.err != nil {
		return b
	}
	b.addStep(&routestep.SetProperty{
		Name:          stepName,
		PropertyName:  propertyName,
		PropertyValue: propertyValue,
	})
	return b
}

func (b *RouteBuilder) To(stepName, uri string) *RouteBuilder {
	if b.err != nil {
		return b
	}
	b.addStep(&routestep.To{
		Name: stepName,
		URI:  uri,
	})
	return b
}

func (b *RouteBuilder) Unmarshal(stepName string, format string, targetType any) *RouteBuilder {
	if b.err != nil {
		return b
	}
	b.addStep(&routestep.Unmarshal{
		Name:       stepName,
		Format:     format,
		TargetType: targetType,
	})
	return b
}

func (b *RouteBuilder) Marshal(stepName string, format string) *RouteBuilder {
	if b.err != nil {
		return b
	}
	b.addStep(&routestep.Marshal{
		Name:   stepName,
		Format: format,
	})
	return b
}

func (b *RouteBuilder) Delay(stepName string, durMs int64) *RouteBuilder {
	if b.err != nil {
		return b
	}
	b.addStep(&routestep.Delay{
		Name:     stepName,
		Duration: durMs,
	})
	return b
}

func (b *RouteBuilder) Log(stepName string, level api.LogLevel, msg string) *RouteBuilder {
	if b.err != nil {
		return b
	}
	b.addStep(&routestep.Log{
		Name:  stepName,
		Msg:   msg,
		Level: level,
	})
	return b
}

func (b *RouteBuilder) LogInfo(stepName, msg string) *RouteBuilder {
	return b.Log(stepName, api.LogLevelInfo, msg)
}

func (b *RouteBuilder) LogWarn(stepName, msg string) *RouteBuilder {
	return b.Log(stepName, api.LogLevelWarn, msg)
}

func (b *RouteBuilder) LogError(stepName, msg string) *RouteBuilder {
	return b.Log(stepName, api.LogLevelError, msg)
}

func (b *RouteBuilder) LogDebug(stepName, msg string) *RouteBuilder {
	return b.Log(stepName, api.LogLevelDebug, msg)
}

// Func adds FuncStep to the current Route.
// userFunc must be:
//  1. string, function must be registered by means of Runtime.RegisterFunc / Runtime.MustRegisterFunc.
//  2. inline fn with type fn(*Exchange).
func (b *RouteBuilder) Func(stepName string, userFunc any) *RouteBuilder {
	if b.err != nil {
		return b
	}

	switch fnt := userFunc.(type) {
	case string, func(*exchange.Exchange):
	default:
		panic(fmt.Errorf("userFunc expected string or fn(*Exchange), but got %T", fnt))
	}

	b.addStep(&routestep.Fn{
		Name: stepName,
		Func: userFunc,
	})
	return b
}

// Pipeline adds new pipeline.
// Function configure will be called to configure PipelineStep.
func (b *RouteBuilder) Pipeline(stepName string, stopOnError bool, configure func(b *RouteBuilder)) *RouteBuilder {
	if b.err != nil {
		return b
	}
	pipe := &routestep.Pipeline{
		Name:       stepName,
		StoOnError: stopOnError,
		Steps:      []api.RouteStep{},
	}
	b.addStep(pipe)

	// push pipeline
	b.pushStack(&pipe.Steps)
	configure(b) // configure PipelineStep
	b.popStack() // pop

	return b
}

func (b *RouteBuilder) Multicast(stepName string) *MulticastStepBuilder {
	if b.err != nil {
		return &MulticastStepBuilder{builder: b}
	}

	multicastStep := &routestep.Multicast{
		Name:        stepName,
		Parallel:    false,
		StopOnError: false,
	}
	b.addStep(multicastStep)

	return &MulticastStepBuilder{builder: b, multicastStep: multicastStep}
}

func (b *RouteBuilder) Loop(stepName string, predicate expr.Expression, copyExchange bool, configure func(b *RouteBuilder)) *RouteBuilder {
	if b.err != nil {
		return b
	}

	step := &routestep.Loop{
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

func (b *RouteBuilder) RemoveHeader(stepName string, headerName ...string) *RouteBuilder {
	if b.err != nil {
		return b
	}
	b.addStep(&routestep.RemoveHeader{
		Name:        stepName,
		HeaderNames: headerName,
	})
	return b
}

func (b *RouteBuilder) RemoveProperty(stepName string, propertyName ...string) *RouteBuilder {
	if b.err != nil {
		return b
	}
	b.addStep(&routestep.RemoveProperty{
		Name:          stepName,
		PropertyNames: propertyName,
	})
	return b
}
