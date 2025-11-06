package routestep

import (
	"fmt"
	"strings"
)

type RemoveProperty struct {
	Name          string
	PropertyNames []string
}

func (s *RemoveProperty) StepName() string {
	if s.Name == "" {
		return fmt.Sprintf("removeProperty[%s]", strings.Join(s.PropertyNames, ","))
	}
	return s.Name
}
