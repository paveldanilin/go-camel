package camel

import (
	"github.com/paveldanilin/go-camel/camel/dsl"
)

// Function aliases from dsl package

type (
	SimpleExpression   func(expression string) dsl.Expression
	ConstantExpression func(value any) dsl.Expression
)

var (
	Simple   SimpleExpression   = dsl.Simple
	Constant ConstantExpression = dsl.Constant
)

func NewRoute(name, from string) *dsl.RouteBuilder {
	return dsl.NewRouteBuilder(name, from)
}
