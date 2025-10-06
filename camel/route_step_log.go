package camel

import "fmt"

type LogStep struct {
	Name string

	// Msg is a string that will be evaluated against *Exchange and send to Runtime logger.
	// Template lang is supported: ${properties.a}
	Msg string

	Level LogLevel
}

func (s *LogStep) StepName() string {
	if s.Name == "" {
		return fmt.Sprintf("log[%d:%s]", s.Level, s.Msg)
	}
	return s.Name
}

// ---------------------------------------------------------------------------------------------------------------------
// LogStepBuilder
// ---------------------------------------------------------------------------------------------------------------------

type LogStepBuilder struct {
	builder *RouteBuilder
	logStep *LogStep
}

func (lb *LogStepBuilder) StepName(stepName string) *LogStepBuilder {
	lb.logStep.Name = stepName
	return lb
}

func (lb *LogStepBuilder) Level(logLevel LogLevel) *LogStepBuilder {
	lb.logStep.Level = logLevel
	return lb
}

func (lb *LogStepBuilder) Msg(msg string) *LogStepBuilder {
	lb.logStep.Msg = msg
	return lb
}

func (lb *LogStepBuilder) EndLog() *RouteBuilder {
	return lb.builder
}

// ---------------------------------------------------------------------------------------------------------------------
// RouteBuilder :: Log
// ---------------------------------------------------------------------------------------------------------------------

// Log adds LogStep to the current Route with the given msg template and LogLevelInfo.
func (b *RouteBuilder) Log(msg string) *RouteBuilder {
	if b.err != nil {
		return b
	}
	b.addStep(&LogStep{
		Msg:   msg,
		Level: LogLevelInfo,
	})
	return b
}

func (b *RouteBuilder) LogWithBuilder() *LogStepBuilder {
	if b.err != nil {
		return &LogStepBuilder{builder: b}
	}

	logStep := &LogStep{}
	b.addStep(logStep)

	return &LogStepBuilder{builder: b, logStep: logStep}
}
