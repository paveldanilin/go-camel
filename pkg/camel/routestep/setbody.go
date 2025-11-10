package routestep

import (
	"fmt"
	"github.com/paveldanilin/go-camel/pkg/camel/expr"
)

type SetBody struct {
	Name      string
	BodyValue expr.Definition
}

func (s *SetBody) StepName() string {
	if s.Name == "" {
		return fmt.Sprintf("setBody[%s:%s]", s.BodyValue.Kind, s.BodyValue.Expression)
	}
	return s.Name
}
