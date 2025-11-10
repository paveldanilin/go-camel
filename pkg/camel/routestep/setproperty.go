package routestep

import (
	"fmt"
	"github.com/paveldanilin/go-camel/pkg/camel/expr"
)

type SetProperty struct {
	Name          string
	PropertyName  string
	PropertyValue expr.Definition
}

func (s *SetProperty) StepName() string {
	if s.Name == "" {
		return fmt.Sprintf("setProperty[%s]={%s:%v}", s.PropertyName, s.PropertyValue.Kind, s.PropertyValue.Expression)
	}
	return s.Name
}
