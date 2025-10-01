package camel

import "fmt"

type SetErrorStep struct {
	Name  string
	Error error
}

func (s *SetErrorStep) StepName() string {
	if s.Name == "" {
		return fmt.Sprintf("setError[%s]", s.Error)
	}
	return s.Name
}

// ---------------------------------------------------------------------------------------------------------------------
// RouteBuilder :: SetError
// ---------------------------------------------------------------------------------------------------------------------

func (b *RouteBuilder) SetError(stepName string, err error) *RouteBuilder {
	if b.err != nil {
		return b
	}
	b.addStep(&SetErrorStep{
		Name:  stepName,
		Error: err,
	})
	return b
}
