package routestep

import (
	"fmt"
	"github.com/paveldanilin/go-camel/pkg/camel/api"
)

type Log struct {
	Name string

	// Msg is a string that will be evaluated against *Exchange and send to Runtime logger.
	// Template lang is supported: ${properties.a}
	Msg string

	Level api.LogLevel
}

func (s *Log) StepName() string {
	if s.Name == "" {
		return fmt.Sprintf("log[%d:%s]", s.Level, s.Msg)
	}
	return s.Name
}
