package step

import (
	"fmt"
	"github.com/paveldanilin/go-camel/pkg/camel/expr"
)

type SetBody struct {
	Name      string
	BodyValue expr.Expression
}

func (s *SetBody) StepName() string {
	if s.Name == "" {
		return fmt.Sprintf("setBody[%s:%s]", s.BodyValue.Language, s.BodyValue.Expression)
	}
	return s.Name
}
