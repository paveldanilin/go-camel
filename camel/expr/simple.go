package expr

import (
	"fmt"
	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
	"github.com/paveldanilin/go-camel/camel"
)

// SimpleExpr is a wrapper for https://expr-lang.org/docs/getting-started
type SimpleExpr struct {
	rawExpr string
	program *vm.Program
}

func Simple(e string) (*SimpleExpr, error) {
	program, err := expr.Compile(e,
		expr.AllowUndefinedVariables(),
		expr.Optimize(true),
		expr.AsAny())
	if err != nil {
		return nil, err
	}

	return &SimpleExpr{
		rawExpr: e,
		program: program,
	}, nil
}

func MustSimple(e string) *SimpleExpr {
	simpleExpr, err := Simple(e)
	if err != nil {
		panic(fmt.Errorf("camel: expr: simple: %w", err))
	}

	return simpleExpr
}

func (e *SimpleExpr) Eval(exchange *camel.Exchange) (any, error) {
	// TODO: move to Exchange?
	env := map[string]any{
		"body":    exchange.Message().Body,
		"headers": exchange.Message().Headers().All(),
		"header":  exchange.Message().Headers().All(),
		"exchange": map[string]any{
			"error":      exchange.Error,
			"properties": exchange.Properties().All(),
		},
	}

	return expr.Run(e.program, env)
}

func (e *SimpleExpr) String() string {
	return e.rawExpr
}
