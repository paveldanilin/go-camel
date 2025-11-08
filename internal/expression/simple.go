package expression

import (
	"fmt"
	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
	"strconv"
)

var simpleExpressionEnv map[string]any

func init() {
	toStringFn := func(i int) string {
		return strconv.Itoa(i)
	}

	simpleExpressionEnv = map[string]any{
		"toString": toStringFn,
	}
}

// simple is a wrapper for https://expr-lang.org/docs/getting-started
// Variables (see Exchange.asMap):
//
//	 body:			the Message body
//	 header:		the Message headers (header.foo refers to the Exchange header 'foo')
//	 error:			the Exchange error
//	 property:		the Exchange properties (property.foo refers to the Exchange property 'foo')
//	 id:				the Message id
//		exchangeId:		the Exchange id
type simple struct {
	raw     string
	program *vm.Program
}

func NewSimple(e string) (*simple, error) {
	program, err := expr.Compile(e,
		expr.AllowUndefinedVariables(),
		expr.Optimize(true),
		expr.AsAny())
	if err != nil {
		return nil, err
	}

	return &simple{
		raw:     e,
		program: program,
	}, nil
}

func MustSimple(e string) *simple {
	expr_, err := NewSimple(e)
	if err != nil {
		panic(fmt.Errorf("camel: expression: simple: %w", err))
	}

	return expr_
}

func (e *simple) Eval(ex *exchange.Exchange) (any, error) {
	env := ex.AsMap()
	for k, v := range simpleExpressionEnv {
		env[k] = v
	}
	return expr.Run(e.program, env)
}

func (e *simple) String() string {
	return e.raw
}
