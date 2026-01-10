package step

import (
	"fmt"
	"github.com/paveldanilin/go-camel/pkg/camel/api"
	"github.com/paveldanilin/go-camel/pkg/camel/route/expr"
)

type WhileLoop struct {
	Name         string
	Predicate    expr.Definition
	CopyExchange bool
	Steps        []api.RouteStep
}

func (s *WhileLoop) setStepName(stepName string) {
	s.Name = stepName
}

func (s *WhileLoop) StepName() string {
	if s.Name == "" {
		return fmt.Sprintf("while[%s:%s]", s.Predicate.Kind, s.Predicate.Expression)
	}
	return s.Name
}
