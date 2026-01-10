package step

import (
	"fmt"
	"github.com/paveldanilin/go-camel/pkg/camel/route/expr"
)

type SetVariable struct {
	Name          string
	VariableName  string
	VariableValue expr.Definition
}

func (s *SetVariable) setStepName(stepName string) {
	s.Name = stepName
}

func (s *SetVariable) StepName() string {
	if s.Name == "" {
		return fmt.Sprintf("setVariable[%s]={%s:%v}", s.VariableName, s.VariableValue.Kind, s.VariableValue.Expression)
	}
	return s.Name
}
