package routestep

import (
	"fmt"
	"github.com/paveldanilin/go-camel/pkg/camel/api"
)

type Pipeline struct {
	Name       string
	StoOnError bool
	Steps      []api.RouteStep
}

func (s *Pipeline) StepName() string {
	if s.Name == "" {
		return fmt.Sprintf("pipeline[stopOnError=%v]", s.StoOnError)
	}
	return s.Name
}
