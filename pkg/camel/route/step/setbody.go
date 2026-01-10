package step

import (
	"fmt"
	"github.com/paveldanilin/go-camel/pkg/camel/route/expr"
)

// SetBody represents a set body route step definition.
type SetBody struct {
	Name      string
	BodyValue expr.Definition
}

func (s *SetBody) setStepName(stepName string) {
	s.Name = stepName
}

func (s *SetBody) StepName() string {
	if s.Name == "" {
		return fmt.Sprintf("setBody[%s:%s]", s.BodyValue.Kind, s.BodyValue.Expression)
	}
	return s.Name
}
