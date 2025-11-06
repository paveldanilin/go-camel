package routestep

import "fmt"

type ConvertBody struct {
	Name       string
	TargetType any
	NamedType  string
	Params     map[string]any
}

func (s *ConvertBody) StepName() string {
	if s.Name == "" {
		if s.TargetType == nil {
			return fmt.Sprintf("convertBody[%s;%v]", s.NamedType, s.Params)
		}
		return fmt.Sprintf("convertBody[%T;%v]", s.TargetType, s.Params)
	}
	return s.Name
}
