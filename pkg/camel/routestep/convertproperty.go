package routestep

import "fmt"

type ConvertProperty struct {
	Name         string
	PropertyName string
	TargetType   any
	NamedType    string
	Params       map[string]any
}

func (s *ConvertProperty) StepName() string {
	if s.Name == "" {
		if s.TargetType == nil {
			return fmt.Sprintf("convertProperty[%s;%s;%v]", s.PropertyName, s.NamedType, s.Params)
		}
		return fmt.Sprintf("convertProperty[%s;%T;%v]", s.PropertyName, s.TargetType, s.Params)
	}
	return s.Name
}
