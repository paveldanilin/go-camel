package expr

import (
	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
	"github.com/paveldanilin/go-camel/camel"
)

type SimpleExpr struct {
	expr string
	prog *vm.Program
}

func Simple(e string) (*SimpleExpr, error) {

	prog, err := expr.Compile(e, expr.AllowUndefinedVariables())
	if err != nil {
		return nil, err
	}

	return &SimpleExpr{
		expr: e,
		prog: prog,
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

	return expr.Run(e.prog, message)
}
