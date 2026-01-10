package step

import (
	"fmt"
	"strings"
)

type RemoveProperty struct {
	Name          string
	PropertyNames []string
}

func (s *RemoveProperty) setStepName(stepName string) {
	s.Name = stepName
}

func (s *RemoveProperty) StepName() string {
	if s.Name == "" {
		return fmt.Sprintf("removeProperty[%s]", strings.Join(s.PropertyNames, ","))
	}
	return s.Name
}
