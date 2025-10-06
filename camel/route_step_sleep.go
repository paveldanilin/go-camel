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

// ---------------------------------------------------------------------------------------------------------------------
// RouteBuilder :: Sleep
// ---------------------------------------------------------------------------------------------------------------------

// Sleep adds SleepStep to the current Route with the given dur (milliseconds).
func (b *RouteBuilder) Sleep(dur int64) *RouteBuilder {
	if b.err != nil {
		return b
	}
	b.addStep(&SleepStep{
		Duration: dur,
	})
	return b
}

// SleepWithName adds SleepStep to the current Route with the given stepName and dur (milliseconds).
func (b *RouteBuilder) SleepWithName(stepName string, dur int64) *RouteBuilder {
	if b.err != nil {
		return b
	}
	b.addStep(&SleepStep{
		Name:     stepName,
		Duration: dur,
	})
	return b
}
