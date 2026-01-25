package step

import (
	"fmt"
	"log/slog"
)

type Log struct {

	// Step name
	Name string

	// Msg is a string that will be evaluated against *Exchange and send to Runtime logger.
	// Supports variable interpolation, like: "Hello, ${property.username}".
	// For more details see: template.Parse.
	Msg string

	// Log level, see slog.Level:
	// 	LevelDebug Level = -4
	//	LevelInfo  Level = 0
	//	LevelWarn  Level = 4
	//	LevelError Level = 8
	Level slog.Level
}

func (s *Log) setStepName(stepName string) {
	s.Name = stepName
}

func (s *Log) StepName() string {
	if s.Name == "" {
		return fmt.Sprintf("log[%d:%s]", s.Level, s.Msg)
	}

	return s.Name
}
