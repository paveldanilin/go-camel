package expr

import (
	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
	"github.com/paveldanilin/go-camel/camel"
)

type SimpleExpr struct {
	expr    string
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
		expr:    e,
		program: program,
	}, nil
}

func MustSimple(e string) *SimpleExpr {

	simpleExpr, err := Simple(e)
	if err != nil {
		panic(err)
	}

	return simpleExpr
}

func (e *SimpleExpr) Eval(message *camel.Message) (any, error) {

	env := map[string]any{
		"body":    message.Body,
		"headers": message.MessageHeaders().Headers(),
		"header":  message.MessageHeaders().Headers(),
		"properties": map[string]any{
			"error": message.Error,
		},
	}

	return expr.Run(e.program, env)
}
