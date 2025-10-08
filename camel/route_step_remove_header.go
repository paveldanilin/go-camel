package camel

import "fmt"

type RemoveHeaderStep struct {
	Name       string
	HeaderName string
}

func (s *RemoveHeaderStep) StepName() string {
	if s.Name == "" {
		return fmt.Sprintf("removeHeader[%s]", s.HeaderName)
	}
	return s.Name
}

// ---------------------------------------------------------------------------------------------------------------------
// RouteBuilder :: RemoveHeader
// ---------------------------------------------------------------------------------------------------------------------

// RemoveHeader adds steps to remove message header.
func (b *RouteBuilder) RemoveHeader(stepName, headerName string) *RouteBuilder {
	if b.err != nil {
		return b
	}
	b.addStep(&RemoveHeaderStep{
		Name:       stepName,
		HeaderName: headerName,
	})
	return b
}
