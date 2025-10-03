package camel

import "fmt"

type SetBodyStep struct {
	Name      string
	BodyValue Expression
}

func (s *SetBodyStep) StepName() string {
	if s.Name == "" {
		return fmt.Sprintf("setBody[%s:%s]", s.BodyValue.Language, s.BodyValue.Expression)
	}
	return s.Name
}

// SetBody adds step to set the message body.
func (b *RouteBuilder) SetBody(stepName string, bodyValue Expression) *RouteBuilder {
	if b.err != nil {
		return b
	}
	b.addStep(&SetBodyStep{
		Name:      stepName,
		BodyValue: bodyValue,
	})
	return b
}
