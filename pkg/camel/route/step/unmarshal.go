package step

import "fmt"

type Unmarshal struct {
	Name       string
	Format     string
	TargetType any
}

func (s *Unmarshal) setStepName(stepName string) {
	s.Name = stepName
}

func (s *Unmarshal) StepName() string {
	if s.Name == "" {
		return fmt.Sprintf("unmarshal[format=%s;targetType=%v]", s.Format, s.TargetType)
	}
	return s.Name
}
