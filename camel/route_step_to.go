package camel

import "fmt"

type ToStep struct {
	Name string
	URI  string
}

func (s *ToStep) StepName() string {
	if s.Name == "" {
		return fmt.Sprintf("to[%s]", s.URI)
	}
	return s.Name
}

// ---------------------------------------------------------------------------------------------------------------------
// RouteBuilder :: To
// ---------------------------------------------------------------------------------------------------------------------

func (b *RouteBuilder) To(stepName, uri string) *RouteBuilder {
	if b.err != nil {
		return b
	}
	b.addStep(&ToStep{
		Name: stepName,
		URI:  uri,
	})
	return b
}
