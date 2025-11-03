package step

import "fmt"

type To struct {
	Name string
	URI  string
}

func (s *To) StepName() string {
	if s.Name == "" {
		return fmt.Sprintf("to[%s]", s.URI)
	}
	return s.Name
}
