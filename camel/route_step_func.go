package camel

import "fmt"

type FuncStep struct {
	Name string

	// string or func(*camel.Exchange)
	Func any
}

func (s *FuncStep) StepName() string {
	if s.Name == "" {
		return fmt.Sprintf("func[%T:%v]", s.Func, s.Name)
	}
	return s.Name
}

// ---------------------------------------------------------------------------------------------------------------------
// RouteBuilder :: Func
// ---------------------------------------------------------------------------------------------------------------------

// Func adds FuncStep to the current Route.
// userFunc must be:
//  1. string, function must be registered by means of Runtime.RegisterFunc / Runtime.MustRegisterFunc.
//  2. inline func with type func(*Exchange).
func (b *RouteBuilder) Func(stepName string, userFunc any) *RouteBuilder {
	if b.err != nil {
		return b
	}

	switch fnt := userFunc.(type) {
	case string, func(*Exchange):
	default:
		panic(fmt.Errorf("userFunc expected string or func(*Exchange), but got %T", fnt))
	}

	b.addStep(&FuncStep{
		Name: stepName,
		Func: userFunc,
	})
	return b
}
