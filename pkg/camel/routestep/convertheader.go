package routestep

import "fmt"

type ConvertHeader struct {
	Name       string
	HeaderName string
	TargetType any
	NamedType  string
	Params     map[string]any
}

func (s *ConvertHeader) StepName() string {
	if s.Name == "" {
		if s.TargetType == nil {
			return fmt.Sprintf("convertHeader[%s;%s;%v]", s.HeaderName, s.NamedType, s.Params)
		}
		return fmt.Sprintf("convertHeader[%s;%T;%v]", s.HeaderName, s.TargetType, s.Params)
	}
	return s.Name
}
