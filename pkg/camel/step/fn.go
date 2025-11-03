package step

import "fmt"

type Fn struct {
	Name string

	// string or fn(*camel.Exchange)
	Func any
}

func (s *Fn) StepName() string {
	if s.Name == "" {
		return fmt.Sprintf("fn[%T:%v]", s.Func, s.Name)
	}
	return s.Name
}
