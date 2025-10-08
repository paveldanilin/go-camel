package camel

import "fmt"

type RemovePropertyStep struct {
	Name         string
	PropertyName string
}

func (s *RemovePropertyStep) StepName() string {
	if s.Name == "" {
		return fmt.Sprintf("removeProperty[%s]", s.PropertyName)
	}
	return s.Name
}

// ---------------------------------------------------------------------------------------------------------------------
// RouteBuilder :: RemoveProperty
// ---------------------------------------------------------------------------------------------------------------------

// RemoveProperty adds step to set an exchange property.
func (b *RouteBuilder) RemoveProperty(stepName, propertyName string) *RouteBuilder {
	if b.err != nil {
		return b
	}
	b.addStep(&RemovePropertyStep{
		Name:         stepName,
		PropertyName: propertyName,
	})
	return b
}
