package route

import (
	"fmt"
	"github.com/paveldanilin/go-camel/pkg/camel/api"
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
	"github.com/paveldanilin/go-camel/pkg/camel/route/expr"
	"github.com/paveldanilin/go-camel/pkg/camel/route/step"
	"strings"
)

// Definition represents camel route definition.
// Route is a set of processing steps that are applied to a message as it travels from a source to a destination.
// A route typically consists of a series of processing steps that are connected in a linear sequence.
type Definition struct {
	Name  string
	From  string
	Steps []api.RouteStep
}

// Builder represents a route Definition builder.
type Builder struct {
	route *Definition
	stack []*[]api.RouteStep
	err   error
}

func NewBuilder(name, from string) *Builder {
	name = strings.TrimSpace(name)
	from = strings.TrimSpace(from)
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
func (b *Builder) SetBody(value expr.Definition, opts ...step.Option) *Builder {
	if b.err != nil {
		return b
	}

	setBody := &step.SetBody{
		BodyValue: value,
	}

	step.ApplyOptions(setBody, opts...)

	b.addStep(setBody)
	return b
}

func (b *Builder) ConvertBodyTo(targetType any, params map[string]any, opts ...step.Option) *Builder {
	if b.err != nil {
		return b
	}

	convertBody := &step.ConvertBody{
		TargetType: targetType,
		Params:     params,
	}

	step.ApplyOptions(convertBody, opts...)

	b.addStep(convertBody)
	return b
}

func (b *Builder) ConvertBodyToNamed(typeName string, params map[string]any, opts ...step.Option) *Builder {
	if b.err != nil {
		return b
	}

	convertBody := &step.ConvertBody{
		NamedType: typeName,
		Params:    params,
	}

	step.ApplyOptions(convertBody, opts...)

	b.addStep(convertBody)
	return b
}

func (b *Builder) ConvertHeaderTo(name string, targetType any, params map[string]any, opts ...step.Option) *Builder {
	if b.err != nil {
		return b
	}

	convertHeader := &step.ConvertHeader{
		HeaderName: name,
		TargetType: targetType,
		Params:     params,
	}

	step.ApplyOptions(convertHeader, opts...)

	b.addStep(convertHeader)
	return b
}

func (b *Builder) ConvertHeaderToNamed(name, typeName string, params map[string]any, opts ...step.Option) *Builder {
	if b.err != nil {
		return b
	}

	convertHeader := &step.ConvertHeader{
		HeaderName: name,
		NamedType:  typeName,
		Params:     params,
	}

	step.ApplyOptions(convertHeader, opts...)

	b.addStep(convertHeader)
	return b
}

func (b *Builder) ConvertPropertyTo(name string, targetType any, params map[string]any, opts ...step.Option) *Builder {
	if b.err != nil {
		return b
	}

	convertProperty := &step.ConvertProperty{
		PropertyName: name,
		TargetType:   targetType,
		Params:       params,
	}

	step.ApplyOptions(convertProperty, opts...)

	b.addStep(convertProperty)
	return b
}

func (b *Builder) ConvertPropertyToNamed(name, typeName string, params map[string]any, opts ...step.Option) *Builder {
	if b.err != nil {
		return b
	}

	convertProperty := &step.ConvertProperty{
		PropertyName: name,
		NamedType:    typeName,
		Params:       params,
	}

	step.ApplyOptions(convertProperty, opts...)

	b.addStep(convertProperty)
	return b
}

func (b *Builder) ConvertVariableTo(name string, targetType any, params map[string]any, opts ...step.Option) *Builder {
	if b.err != nil {
		return b
	}

	convertVariable := &step.ConvertVariable{
		VariableName: name,
		TargetType:   targetType,
		Params:       params,
	}

	step.ApplyOptions(convertVariable, opts...)

	b.addStep(convertVariable)
	return b
}

func (b *Builder) ConvertVariableToNamed(name, typeName string, params map[string]any, opts ...step.Option) *Builder {
	if b.err != nil {
		return b
	}

	convertVariable := &step.ConvertVariable{
		VariableName: name,
		NamedType:    typeName,
		Params:       params,
	}

	step.ApplyOptions(convertVariable, opts...)

	b.addStep(convertVariable)
	return b
}

// Choice adds choice step to the current step level.
func (b *Builder) Choice() *ChoiceBuilder {
	if b.err != nil {
		return &ChoiceBuilder{builder: b}
	}

	choice := &step.Choice{}

	b.addStep(choice)

	return &ChoiceBuilder{builder: b, choiceStep: choice}
}

// SetHeader adds steps to set message header.
func (b *Builder) SetHeader(name string, value expr.Definition, opts ...step.Option) *Builder {
	if b.err != nil {
		return b
	}

	setHeader := &step.SetHeader{
		HeaderName:  name,
		HeaderValue: value,
	}

	step.ApplyOptions(setHeader, opts...)

	b.addStep(setHeader)
	return b
}

// Try adds try-catch-finally step.
func (b *Builder) Try(block func(b *Builder)) *TryBuilder {
	if b.err != nil {
		return &TryBuilder{builder: b}
	}

	tryStep := &step.Try{}
	b.addStep(tryStep)

	b.pushStack(&tryStep.Steps)
	block(b) // block try
	b.popStack()

	return &TryBuilder{builder: b, tryStep: tryStep}
}

func (b *Builder) SetError(err error, opts ...step.Option) *Builder {
	if b.err != nil {
		return b
	}

	setError := &step.SetError{
		Error: err,
	}

	step.ApplyOptions(setError, opts...)

	b.addStep(setError)
	return b
}

// SetProperty adds step to set an exchange property.
func (b *Builder) SetProperty(name string, value expr.Definition, opts ...step.Option) *Builder {
	if b.err != nil {
		return b
	}

	setProperty := &step.SetProperty{
		PropertyName:  name,
		PropertyValue: value,
	}

	step.ApplyOptions(setProperty, opts...)

	b.addStep(setProperty)
	return b
}

func (b *Builder) SetVariable(name string, value expr.Definition, opts ...step.Option) *Builder {
	if b.err != nil {
		return b
	}

	setVariable := &step.SetVariable{
		VariableName:  name,
		VariableValue: value,
	}

	step.ApplyOptions(setVariable, opts...)

	b.addStep(setVariable)
	return b
}

func (b *Builder) To(uri string, opts ...step.Option) *Builder {
	if b.err != nil {
		return b
	}

	to := &step.To{
		URI: uri,
	}

	step.ApplyOptions(to, opts...)

	b.addStep(to)
	return b
}

func (b *Builder) Unmarshal(format string, targetType any, opts ...step.Option) *Builder {
	if b.err != nil {
		return b
	}

	unmarshal := &step.Unmarshal{
		Format:     format,
		TargetType: targetType,
	}

	step.ApplyOptions(unmarshal, opts...)

	b.addStep(unmarshal)
	return b
}

func (b *Builder) Marshal(format string, opts ...step.Option) *Builder {
	if b.err != nil {
		return b
	}

	marshal := &step.Marshal{
		Format: format,
	}

	step.ApplyOptions(marshal, opts...)

	b.addStep(marshal)
	return b
}

func (b *Builder) Delay(durMs int64, opts ...step.Option) *Builder {
	if b.err != nil {
		return b
	}

	delay := &step.Delay{
		Duration: durMs,
	}

	step.ApplyOptions(delay, opts...)

	b.addStep(delay)
	return b
}

func (b *Builder) Log(level api.LogLevel, msg string, opts ...step.Option) *Builder {
	if b.err != nil {
		return b
	}

	log := &step.Log{
		Msg:   msg,
		Level: level,
	}

	step.ApplyOptions(log, opts...)

	b.addStep(log)
	return b
}

func (b *Builder) LogInfo(msg string, opts ...step.Option) *Builder {
	return b.Log(api.LogLevelInfo, msg, opts...)
}

func (b *Builder) LogWarn(msg string, opts ...step.Option) *Builder {
	return b.Log(api.LogLevelWarn, msg, opts...)
}

func (b *Builder) LogError(msg string, opts ...step.Option) *Builder {
	return b.Log(api.LogLevelError, msg, opts...)
}

func (b *Builder) LogDebug(msg string, opts ...step.Option) *Builder {
	return b.Log(api.LogLevelDebug, msg, opts...)
}

// Func adds step.Fn to the current Route.
// userFunc must be one of:
//  1. string, function must be registered by means of Runtime.RegisterFunc / Runtime.MustRegisterFunc.
//  2. inline fn with type fn(*Exchange).
func (b *Builder) Func(userFunc any, opts ...step.Option) *Builder {
	if b.err != nil {
		return b
	}

	switch fnt := userFunc.(type) {
	case string, func(*exchange.Exchange):
	default:
		panic(fmt.Errorf("userFunc expected string or fn(*Exchange), but got %T", fnt))
	}

	fn := &step.Fn{
		Func: userFunc,
	}

	step.ApplyOptions(fn, opts...)

	b.addStep(fn)
	return b
}

// Pipeline adds new pipeline.
// Function configure will be called to configure PipelineStep.
/*func (b *Builder) Pipeline(stepName string, stopOnError bool, configure func(b *Builder)) *Builder {
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
}*/

func (b *Builder) Multicast() *MulticastBuilder {
	if b.err != nil {
		return &MulticastBuilder{builder: b}
	}

	multicastStep := &step.Multicast{
		Parallel:    false,
		StopOnError: false,
	}
	b.addStep(multicastStep)

	return &MulticastBuilder{builder: b, multicastStep: multicastStep}
}

func (b *Builder) While(predicate expr.Definition) *WhileBuilder {
	if b.err != nil {
		return &WhileBuilder{builder: b}
	}

	loopStep := &step.WhileLoop{
		Predicate: predicate,
	}
	b.addStep(loopStep)

	return &WhileBuilder{builder: b, loopStep: loopStep}
}

func (b *Builder) RemoveHeader(names []string, opts ...step.Option) *Builder {
	if b.err != nil {
		return b
	}

	removeHeader := &step.RemoveHeader{
		HeaderNames: names,
	}

	step.ApplyOptions(removeHeader, opts...)

	b.addStep(removeHeader)
	return b
}

func (b *Builder) RemoveProperty(names []string, opts ...step.Option) *Builder {
	if b.err != nil {
		return b
	}

	removeProperty := &step.RemoveProperty{
		PropertyNames: names,
	}

	step.ApplyOptions(removeProperty, opts...)

	b.addStep(removeProperty)
	return b
}

func (b *Builder) RemoveVariable(names []string, opts ...step.Option) *Builder {
	if b.err != nil {
		return b
	}

	removeVariable := &step.RemoveVariable{
		VariableNames: names,
	}

	step.ApplyOptions(removeVariable, opts...)

	b.addStep(removeVariable)
	return b
}
