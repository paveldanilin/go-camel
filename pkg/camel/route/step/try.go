package step

import (
	"fmt"
	"github.com/paveldanilin/go-camel/pkg/camel/route/errmatcher"
	"strings"
)

type CatchWhen struct {
	ErrorMatcher errmatcher.Definition
	Steps        []RouteStep
}

type Try struct {
	Name         string
	Steps        []RouteStep
	WhenCatches  []CatchWhen
	FinallySteps []RouteStep
}

func (s *Try) setStepName(stepName string) {
	s.Name = stepName
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
