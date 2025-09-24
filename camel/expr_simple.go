package camel

import (
	"fmt"
	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
)

// simpleExpr is a wrapper for https://expr-lang.org/docs/getting-started
type simpleExpr struct {
	rawExpr string
	program *vm.Program
}

func newSimpleExpr(e string) (*simpleExpr, error) {
	program, err := expr.Compile(e,
		expr.AllowUndefinedVariables(),
		expr.Optimize(true),
		expr.AsAny())
	if err != nil {
		return nil, err
	}

	return &simpleExpr{
		rawExpr: e,
		program: program,
	}, nil
}

func mustSimpleExpr(e string) *simpleExpr {
	simpleExpr, err := newSimpleExpr(e)
	if err != nil {
		panic(fmt.Errorf("camel: expr: simple: %w", err))
	}

	return simpleExpr
}

func (e *simpleExpr) Eval(exchange *Exchange) (any, error) {
	// TODO: move to Exchange?
	env := map[string]any{
		"body":    exchange.Message().Body,
		"headers": exchange.Message().Headers().All(),
		"header":  exchange.Message().Headers().All(),
		"error":   exchange.Error(),
		"exchange": map[string]any{
			"error":      exchange.Error(),
			"properties": exchange.Properties().All(),
		},
	}

	return expr.Run(e.program, env)
}

func (e *simpleExpr) String() string {
	return e.rawExpr
}
