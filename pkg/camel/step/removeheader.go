package step

import "fmt"

type RemoveHeader struct {
	Name       string
	HeaderName string
}

func (s *RemoveHeader) StepName() string {
	if s.Name == "" {
		return fmt.Sprintf("removeHeader[%s]", s.HeaderName)
	}
	return s.Name
}
