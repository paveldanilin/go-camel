package step

import "fmt"

type SetError struct {
	Name  string
	Error error
}

func (s *SetError) StepName() string {
	if s.Name == "" {
		return fmt.Sprintf("setError[%s]", s.Error)
	}
	return s.Name
}
