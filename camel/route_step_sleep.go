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
func (b *RouteBuilder) Sleep(stepName string, durMs int64) *RouteBuilder {
	if b.err != nil {
		return b
	}
	b.addStep(&SleepStep{
		Name:     stepName,
		Duration: durMs,
	})
	return b
}
