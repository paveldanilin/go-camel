package routestep

import (
	"fmt"
	"github.com/paveldanilin/go-camel/pkg/camel/api"
	"strings"
)

type OutputProcess struct {
	Steps []api.RouteStep
}

func (out OutputProcess) StepName() string {
	return "output"
}

type Multicast struct {
	Name        string
	Parallel    bool
	StopOnError bool
	Aggregator  api.ExchangeAggregator
	Outputs     []OutputProcess
}

func (s *Multicast) StepName() string {
	if s.Name == "" {
		outputNames := make([]string, len(s.Outputs))
		for i, output := range s.Outputs {
			if len(output.Steps) > 0 {
				outputNames[i] = output.Steps[0].StepName()
			} else {
				outputNames[i] = fmt.Sprintf("%d", i)
			}
		}
		return fmt.Sprintf("multicast[parallel=%v;stopOnError=%v;outputs=%s]",
			s.Parallel, s.StopOnError, strings.Join(outputNames, ";"))
	}
	return s.Name
}
