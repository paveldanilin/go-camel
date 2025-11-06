package routestep

import (
	"fmt"
	"strings"
)

type RemoveHeader struct {
	Name        string
	HeaderNames []string
}

func (s *RemoveHeader) StepName() string {
	if s.Name == "" {
		return fmt.Sprintf("removeHeader[%s]", strings.Join(s.HeaderNames, ","))
	}
	return s.Name
}
