package step

import "fmt"

type ConvertVariable struct {
	Name         string
	VariableName string
	TargetType   any
	NamedType    string
	Params       map[string]any
}

func (s *ConvertVariable) setStepName(stepName string) {
	s.Name = stepName
}

func (s *ConvertVariable) StepName() string {
	if s.Name == "" {
		if s.TargetType == nil {
			return fmt.Sprintf("convertVariable[%s;%s;%v]", s.VariableName, s.NamedType, s.Params)
		}
		return fmt.Sprintf("convertVariable[%s;%T;%v]", s.VariableName, s.TargetType, s.Params)
	}
	return s.Name
}
