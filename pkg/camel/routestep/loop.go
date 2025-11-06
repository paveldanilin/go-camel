package routestep

import (
	"fmt"
	"github.com/paveldanilin/go-camel/pkg/camel/api"
	"github.com/paveldanilin/go-camel/pkg/camel/expr"
)

type Loop struct {
	Name         string
	Predicate    expr.Expression
	CopyExchange bool
	Steps        []api.RouteStep
}

func (s *Loop) StepName() string {
	if s.Name == "" {
		return fmt.Sprintf("loop[%s:%s]", s.Predicate.Language, s.Predicate.Expression)
	}
	return s.Name
}
