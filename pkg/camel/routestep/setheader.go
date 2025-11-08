package routestep

import (
	"fmt"
	"github.com/paveldanilin/go-camel/pkg/camel/expr"
)

type SetHeader struct {
	Name        string
	HeaderName  string
	HeaderValue expr.Definition
}

func (s *SetHeader) StepName() string {
	if s.Name == "" {
		return fmt.Sprintf("setHeader[%s]={%s:%v}", s.HeaderName, s.HeaderValue.Kind, s.HeaderValue.Expression)
	}
	return s.Name
}
