package camel

import "fmt"

type MarshalStep struct {
	Name   string
	Format string
}

func (s *MarshalStep) StepName() string {
	if s.Name == "" {
		return fmt.Sprintf("marshal[format=%s]", s.Format)
	}
	return s.Name
}

// ---------------------------------------------------------------------------------------------------------------------
// RouteBuilder :: Marshal
// ---------------------------------------------------------------------------------------------------------------------

func (b *RouteBuilder) Marshal(stepName string, format string) *RouteBuilder {
	if b.err != nil {
		return b
	}
	b.addStep(&MarshalStep{
		Name:   stepName,
		Format: format,
	})
	return b
}
