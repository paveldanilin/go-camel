package step

import (
	"fmt"
	"strings"
)

type RemoveVariable struct {
	Name          string
	VariableNames []string
}

func (s *RemoveVariable) setStepName(stepName string) {
	s.Name = stepName
}

func (s *RemoveVariable) StepName() string {
	if s.Name == "" {
		return fmt.Sprintf("removeVariable[%s]", strings.Join(s.VariableNames, ","))
	}
	return s.Name
}
