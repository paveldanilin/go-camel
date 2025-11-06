package routestep

import "fmt"

type Delay struct {
	Name     string
	Duration int64 // milliseconds
}

func (s *Delay) StepName() string {
	if s.Name == "" {
		return fmt.Sprintf("delay[%dms]", s.Duration)
	}
	return s.Name
}
