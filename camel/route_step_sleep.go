package camel

import "fmt"

type SleepStep struct {
	Name     string
	Duration int64 // milliseconds
}

func (s *SleepStep) StepName() string {
	if s.Name == "" {
		return fmt.Sprintf("sleep[%dms]", s.Duration)
	}
	return s.Name
}

func (b *RouteBuilder) Sleep(stepName string, dur int64) *RouteBuilder {
	if b.err != nil {
		return b
	}
	b.addStep(&SleepStep{
		Name:     stepName,
		Duration: dur,
	})
	return b
}
