package dsl

import "fmt"

type FuncStep struct {
	Name string
	Func any
}

func (s *FuncStep) StepName() string {
	if s.Name == "" {
		return fmt.Sprintf("func[%T:%v]", s.Func, s.Name)
	}
	return s.Name
}

func (b *RouteBuilder) Func(stepName string, userFunc any) *RouteBuilder {
	if b.err != nil {
		return b
	}
	b.addStep(&FuncStep{
		Name: stepName,
		Func: userFunc,
	})
	return b
}
