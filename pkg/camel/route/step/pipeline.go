package step

import (
	"fmt"
)

type Pipeline struct {
	Name       string
	StoOnError bool
	Steps      []RouteStep
}

func (s *Pipeline) setStepName(stepName string) {
	s.Name = stepName
}

func (s *Pipeline) StepName() string {
	if s.Name == "" {
		return fmt.Sprintf("pipeline[stopOnError=%v]", s.StoOnError)
	}
	return s.Name
}
