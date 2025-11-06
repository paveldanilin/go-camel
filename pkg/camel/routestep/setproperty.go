package routestep

import (
	"fmt"
	"github.com/paveldanilin/go-camel/pkg/camel/expr"
)

type SetProperty struct {
	Name          string
	PropertyName  string
	PropertyValue expr.Expression
}

func (s *SetProperty) StepName() string {
	if s.Name == "" {
		return fmt.Sprintf("setProperty[%s]={%s:%v}", s.PropertyName, s.PropertyValue.Language, s.PropertyValue.Expression)
	}
	return s.Name
}
