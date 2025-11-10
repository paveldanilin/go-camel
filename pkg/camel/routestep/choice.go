package routestep

import (
	"fmt"
	"github.com/paveldanilin/go-camel/pkg/camel/api"
	"github.com/paveldanilin/go-camel/pkg/camel/expr"
	"strings"
)

type ChoiceWhen struct {
	Predicate expr.Definition
	Steps     []api.RouteStep
}

func (s *ChoiceWhen) StepName() string {
	return fmt.Sprintf("choiceWhen[%s:%s]", s.Predicate.Kind, s.Predicate.Expression)
}

type Choice struct {
	Name      string
	WhenCases []ChoiceWhen
	Otherwise []api.RouteStep
}

func (s *Choice) StepName() string {
	if s.Name == "" {
		when := make([]string, len(s.WhenCases))
		for i, w := range s.WhenCases {
			when[i] = w.StepName()
		}
		return fmt.Sprintf("choice[%s]", strings.Join(when, ";"))
	}
	return s.Name
}
