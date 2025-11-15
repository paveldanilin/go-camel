package route

import (
	"fmt"
	"github.com/paveldanilin/go-camel/pkg/camel/api"
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
	"github.com/paveldanilin/go-camel/pkg/camel/route/expr"
	"github.com/paveldanilin/go-camel/pkg/camel/route/step"
)

// Definition represents camel route definition.
// Route is a set of processing steps that are applied to a message as it travels from a source to a destination.
// A route typically consists of a series of processing steps that are connected in a linear sequence.
type Definition struct {
	Name  string
	From  string
	Steps []api.RouteStep
}

// Builder represents a Route builder.
type Builder struct {
	route *Definition
	stack []*[]api.RouteStep
	err   error
}

func NewBuilder(name, from string) *Builder {
	if name == "" {
		panic(fmt.Errorf("'name' must be not empty string"))
	}
	if from == "" {
		panic(fmt.Errorf("'from' must be not empty string"))
	}
	r := &Definition{
		Name:  name,
		From:  from,
		Steps: []api.RouteStep{},
	}
	return &Builder{
		route: r,
		stack: []*[]api.RouteStep{&r.Steps},
	}
}

func (b *Builder) Build() (*Definition, error) {
	if b.err != nil {
		return nil, b.err
	}

	// Check size != 1 then not all steps where closed properly.
	if len(b.stack) != 1 {
		return nil, fmt.Errorf("not all nested steps where closed properly (PipelineStep, ChoiceStep,... etc)")
	}
	return b.route, nil
}

// currentSteps returns the top of the stack.
func (b *Builder) currentSteps() *[]api.RouteStep {
	return b.stack[len(b.stack)-1]
}

// addStep adds step to the current context.
func (b *Builder) addStep(step api.RouteStep) {
	if b.err != nil {
		return
	}
	steps := b.currentSteps()
	*steps = append(*steps, step)
}

// pushStack pushes builder to the new context.
func (b *Builder) pushStack(steps *[]api.RouteStep) {
	b.stack = append(b.stack, steps)
}

// popStack pops builder.
func (b *Builder) popStack() {
	if len(b.stack) <= 1 {
		b.err = fmt.Errorf("camel: step builder: you are trying to close more block that were opened")
		return
	}
	b.stack = b.stack[:len(b.stack)-1]
}

// SetBody adds step to set the message body.
func (b *Builder) SetBody(stepName string, bodyValue expr.Definition) *Builder {
	if b.err != nil {
		return b
	}
	b.addStep(&step.SetBody{
		Name:      stepName,
		BodyValue: bodyValue,
	})
	return b
}

func (b *Builder) ConvertBodyTo(stepName string, targetType any, params map[string]any) *Builder {
	if b.err != nil {
		return b
	}
	b.addStep(&step.ConvertBody{
		Name:       stepName,
		TargetType: targetType,
		Params:     params,
	})
	return b
}

func (b *Builder) ConvertBodyToNamed(stepName, typeName string, params map[string]any) *Builder {
	if b.err != nil {
		return b
	}
	b.addStep(&step.ConvertBody{
		Name:      stepName,
		NamedType: typeName,
		Params:    params,
	})
	return b
}

func (b *Builder) ConvertHeaderTo(stepName, headerName string, targetType any, params map[string]any) *Builder {
	if b.err != nil {
		return b
	}
	b.addStep(&step.ConvertHeader{
		Name:       stepName,
		TargetType: targetType,
		HeaderName: headerName,
		Params:     params,
	})
	return b
}

func (b *Builder) ConvertHeaderToNamed(stepName, headerName, typeName string, params map[string]any) *Builder {
	if b.err != nil {
		return b
	}
	b.addStep(&step.ConvertHeader{
		Name:       stepName,
		NamedType:  typeName,
		HeaderName: headerName,
		Params:     params,
	})
	return b
}

func (b *Builder) ConvertPropertyTo(stepName, propertyName string, targetType any, params map[string]any) *Builder {
	if b.err != nil {
		return b
	}
	b.addStep(&step.ConvertProperty{
		Name:         stepName,
		TargetType:   targetType,
		PropertyName: propertyName,
		Params:       params,
	})
	return b
}

func (b *Builder) ConvertPropertyToNamed(stepName, propertyName, typeName string, params map[string]any) *Builder {
	if b.err != nil {
		return b
	}
	b.addStep(&step.ConvertProperty{
		Name:         stepName,
		NamedType:    typeName,
		PropertyName: propertyName,
		Params:       params,
	})
	return b
}

// Choice adds choice step to the current step level.
func (b *Builder) Choice(stepName string) *ChoiceStepBuilder {
	if b.err != nil {
		return &ChoiceStepBuilder{builder: b}
	}

	choiceStep := &step.Choice{Name: stepName}
	b.addStep(choiceStep)

	return &ChoiceStepBuilder{builder: b, choiceStep: choiceStep}
}

// SetHeader adds steps to set message header.
func (b *Builder) SetHeader(stepName, headerName string, headerValue expr.Definition) *Builder {
	if b.err != nil {
		return b
	}
	b.addStep(&step.SetHeader{
		Name:        stepName,
		HeaderName:  headerName,
		HeaderValue: headerValue,
	})
	return b
}

// Try adds try-catch-finally step.
func (b *Builder) Try(stepName string, configure func(b *Builder)) *TryStepBuilder {
	if b.err != nil {
		return &TryStepBuilder{builder: b}
	}

	tryStep := &step.Try{Name: stepName}
	b.addStep(tryStep)

	b.pushStack(&tryStep.Steps)
	configure(b) // configure try
	b.popStack()

	return &TryStepBuilder{builder: b, tryStep: tryStep}
}

func (b *Builder) SetError(stepName string, err error) *Builder {
	if b.err != nil {
		return b
	}
	b.addStep(&step.SetError{
		Name:  stepName,
		Error: err,
	})
	return b
}

// SetProperty adds step to set an exchange property.
func (b *Builder) SetProperty(stepName, propertyName string, propertyValue expr.Definition) *Builder {
	if b.err != nil {
		return b
	}
	b.addStep(&step.SetProperty{
		Name:          stepName,
		PropertyName:  propertyName,
		PropertyValue: propertyValue,
	})
	return b
}

func (b *Builder) To(stepName, uri string) *Builder {
	if b.err != nil {
		return b
	}
	b.addStep(&step.To{
		Name: stepName,
		URI:  uri,
	})
	return b
}

func (b *Builder) Unmarshal(stepName string, format string, targetType any) *Builder {
	if b.err != nil {
		return b
	}
	b.addStep(&step.Unmarshal{
		Name:       stepName,
		Format:     format,
		TargetType: targetType,
	})
	return b
}

func (b *Builder) Marshal(stepName string, format string) *Builder {
	if b.err != nil {
		return b
	}
	b.addStep(&step.Marshal{
		Name:   stepName,
		Format: format,
	})
	return b
}

func (b *Builder) Delay(stepName string, durMs int64) *Builder {
	if b.err != nil {
		return b
	}
	b.addStep(&step.Delay{
		Name:     stepName,
		Duration: durMs,
	})
	return b
}

func (b *Builder) Log(stepName string, level api.LogLevel, msg string) *Builder {
	if b.err != nil {
		return b
	}
	b.addStep(&step.Log{
		Name:  stepName,
		Msg:   msg,
		Level: level,
	})
	return b
}

func (b *Builder) LogInfo(stepName, msg string) *Builder {
	return b.Log(stepName, api.LogLevelInfo, msg)
}

func (b *Builder) LogWarn(stepName, msg string) *Builder {
	return b.Log(stepName, api.LogLevelWarn, msg)
}

func (b *Builder) LogError(stepName, msg string) *Builder {
	return b.Log(stepName, api.LogLevelError, msg)
}

func (b *Builder) LogDebug(stepName, msg string) *Builder {
	return b.Log(stepName, api.LogLevelDebug, msg)
}

// Func adds FuncStep to the current Route.
// userFunc must be:
//  1. string, function must be registered by means of Runtime.RegisterFunc / Runtime.MustRegisterFunc.
//  2. inline fn with type fn(*Exchange).
func (b *Builder) Func(stepName string, userFunc any) *Builder {
	if b.err != nil {
		return b
	}

	switch fnt := userFunc.(type) {
	case string, func(*exchange.Exchange):
	default:
		panic(fmt.Errorf("userFunc expected string or fn(*Exchange), but got %T", fnt))
	}

	b.addStep(&step.Fn{
		Name: stepName,
		Func: userFunc,
	})
	return b
}

// Pipeline adds new pipeline.
// Function configure will be called to configure PipelineStep.
func (b *Builder) Pipeline(stepName string, stopOnError bool, configure func(b *Builder)) *Builder {
	if b.err != nil {
		return b
	}
	pipe := &step.Pipeline{
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

func (b *Builder) Multicast(stepName string) *MulticastStepBuilder {
	if b.err != nil {
		return &MulticastStepBuilder{builder: b}
	}

	multicastStep := &step.Multicast{
		Name:        stepName,
		Parallel:    false,
		StopOnError: false,
	}
	b.addStep(multicastStep)

	return &MulticastStepBuilder{builder: b, multicastStep: multicastStep}
}

func (b *Builder) Loop(stepName string, predicate expr.Definition, copyExchange bool, configure func(b *Builder)) *Builder {
	if b.err != nil {
		return b
	}

	step := &step.Loop{
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

func (b *Builder) RemoveHeader(stepName string, headerName ...string) *Builder {
	if b.err != nil {
		return b
	}
	b.addStep(&step.RemoveHeader{
		Name:        stepName,
		HeaderNames: headerName,
	})
	return b
}

func (b *Builder) RemoveProperty(stepName string, propertyName ...string) *Builder {
	if b.err != nil {
		return b
	}
	b.addStep(&step.RemoveProperty{
		Name:          stepName,
		PropertyNames: propertyName,
	})
	return b
}
