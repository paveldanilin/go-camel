package expression

import (
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
)

// Expression represents an expression that takes Exchange and returns computed valueExpression or error.
// Used to dynamically compute valueExpression to setBodyProcessor/setHeaderProcessor.
type Expression interface {
	Eval(e *exchange.Exchange) (any, error)
}
