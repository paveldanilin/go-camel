package camel

import "fmt"

type UnmarshalStep struct {
	Name       string
	Format     string
	TargetType any
}

func (s *UnmarshalStep) StepName() string {
	if s.Name == "" {
		return fmt.Sprintf("unmarshal[format=%s;targetType=%v]", s.Format, s.TargetType)
	}
	return s.Name
}

// ---------------------------------------------------------------------------------------------------------------------
// RouteBuilder :: Unmarshal
// ---------------------------------------------------------------------------------------------------------------------

func (b *RouteBuilder) Unmarshal(stepName string, format string, targetType any) *RouteBuilder {
	if b.err != nil {
		return b
	}
	b.addStep(&UnmarshalStep{
		Name:       stepName,
		Format:     format,
		TargetType: targetType,
	})
	return b
}
