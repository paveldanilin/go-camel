package camel

import (
	"github.com/paveldanilin/go-camel/dsl"
)

// Function aliases from dsl package

type (
	SimpleExpressionFunc   func(expression string) dsl.Expression
	ConstantExpressionFunc func(value any) dsl.Expression
	ErrAnyFunc             func() dsl.ErrorMatcher
	ErrEqualsFunc          func(str string) dsl.ErrorMatcher
	ErrIsFunc              func(target error) dsl.ErrorMatcher
	ErrContainsFunc        func(str string) dsl.ErrorMatcher
	ErrMatchesFunc         func(pattern string) dsl.ErrorMatcher
)

var (
	Simple      SimpleExpressionFunc   = dsl.Simple
	Constant    ConstantExpressionFunc = dsl.Constant
	ErrAny      ErrAnyFunc             = dsl.ErrAny
	ErrEquals   ErrEqualsFunc          = dsl.ErrEquals
	ErrIs       ErrIsFunc              = dsl.ErrIs
	ErrContains ErrContainsFunc        = dsl.ErrContains
	ErrMatches  ErrMatchesFunc         = dsl.ErrMatches
)

func NewRoute(name, from string) *dsl.RouteBuilder {
	return dsl.NewRouteBuilder(name, from)
}
