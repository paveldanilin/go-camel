package routestep

import "fmt"

type Marshal struct {
	Name   string
	Format string
}

func (s *Marshal) StepName() string {
	if s.Name == "" {
		return fmt.Sprintf("marshal[format=%s]", s.Format)
	}
	return s.Name
}
