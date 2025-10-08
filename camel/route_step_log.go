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
// RouteBuilder :: Log
// ---------------------------------------------------------------------------------------------------------------------

// Log adds LogStep to the current Route with the given msg template and LogLevelInfo.
func (b *RouteBuilder) Log(stepName string, level LogLevel, msg string) *RouteBuilder {
	if b.err != nil {
		return b
	}
	b.addStep(&LogStep{
		Name:  stepName,
		Msg:   msg,
		Level: level,
	})
	return b
}

func (b *RouteBuilder) LogInfo(stepName, msg string) *RouteBuilder {
	return b.Log(stepName, LogLevelInfo, msg)
}

func (b *RouteBuilder) LogWarn(stepName, msg string) *RouteBuilder {
	return b.Log(stepName, LogLevelWarn, msg)
}

func (b *RouteBuilder) LogError(stepName, msg string) *RouteBuilder {
	return b.Log(stepName, LogLevelError, msg)
}

func (b *RouteBuilder) LogDebug(stepName, msg string) *RouteBuilder {
	return b.Log(stepName, LogLevelDebug, msg)
}
