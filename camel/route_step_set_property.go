package camel

import "fmt"

type SetPropertyStep struct {
	Name          string
	PropertyName  string
	PropertyValue Expression
}

func (s *SetPropertyStep) StepName() string {
	if s.Name == "" {
		return fmt.Sprintf("setProperty[%s]={%s:%v}", s.PropertyName, s.PropertyValue.Language, s.PropertyValue.Expression)
	}
	return s.Name
}

// ---------------------------------------------------------------------------------------------------------------------
// RouteBuilder :: SetProperty
// ---------------------------------------------------------------------------------------------------------------------

// SetProperty adds step to set an exchange property.
func (b *RouteBuilder) SetProperty(stepName, propertyName string, propertyValue Expression) *RouteBuilder {
	if b.err != nil {
		return b
	}
	b.addStep(&SetPropertyStep{
		Name:          stepName,
		PropertyName:  propertyName,
		PropertyValue: propertyValue,
	})
	return b
}
