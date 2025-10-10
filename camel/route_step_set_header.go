package camel

import "fmt"

type SetHeaderStep struct {
	Name        string
	HeaderName  string
	HeaderValue Expression
}

func (s *SetHeaderStep) StepName() string {
	if s.Name == "" {
		return fmt.Sprintf("setHeader[%s]={%s:%v}", s.HeaderName, s.HeaderValue.Language, s.HeaderValue.Expression)
	}
	return s.Name
}

// ---------------------------------------------------------------------------------------------------------------------
// RouteBuilder :: SetHeader
// ---------------------------------------------------------------------------------------------------------------------

// SetHeader adds steps to set message header.
func (b *RouteBuilder) SetHeader(stepName, headerName string, headerValue Expression) *RouteBuilder {
	if b.err != nil {
		return b
	}
	b.addStep(&SetHeaderStep{
		Name:        stepName,
		HeaderName:  headerName,
		HeaderValue: headerValue,
	})
	return b
}
