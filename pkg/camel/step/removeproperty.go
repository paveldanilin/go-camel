package step

import "fmt"

type RemoveProperty struct {
	Name         string
	PropertyName string
}

func (s *RemoveProperty) StepName() string {
	if s.Name == "" {
		return fmt.Sprintf("removeProperty[%s]", s.PropertyName)
	}
	return s.Name
}
