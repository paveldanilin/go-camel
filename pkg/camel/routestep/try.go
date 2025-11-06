package routestep

import (
	"fmt"
	"github.com/paveldanilin/go-camel/pkg/camel/api"
	"github.com/paveldanilin/go-camel/pkg/camel/errs"
	"strings"
)

type CatchWhen struct {
	ErrorMatcher errs.Matcher
	Steps        []api.RouteStep
}

type Try struct {
	Name         string
	Steps        []api.RouteStep
	WhenCatches  []CatchWhen
	FinallySteps []api.RouteStep
}

func (s *Try) StepName() string {
	if s.Name == "" {
		when := make([]string, len(s.WhenCatches))
		for i, w := range s.WhenCatches {
			when[i] = fmt.Sprintf("%s:%s", w.ErrorMatcher.MatchMode, w.ErrorMatcher.Target)
		}
		return fmt.Sprintf("try[%s]", strings.Join(when, ";"))
	}
	return s.Name
}
