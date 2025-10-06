package camel

import (
	"fmt"
	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
	"strconv"
)

// expression represents an expression that takes Exchange and returns computed valueExpression or error.
// Used to dynamically compute valueExpression to setBodyProcessor/setHeaderProcessor.
type expression interface {
	eval(exchange *Exchange) (any, error)
}

// constExpr represents a constant value.
type constExpr struct {
	value any
}

func newConstExpr(value any) *constExpr {
	return &constExpr{
		value: value,
	}
}

func (e *constExpr) eval(_ *Exchange) (any, error) {
	return e.value, nil
}

// funcExpr represents a user function that returns valueExpression.
type funcExpr func(exchange *Exchange) (any, error)

func (fn funcExpr) eval(exchange *Exchange) (any, error) {
	return fn(exchange)
}

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
	expr_, err := newSimpleExpr(e)
	if err != nil {
		panic(fmt.Errorf("camel: expr: simple: %w", err))
	}

	return expr_
}

func (e *simpleExpr) eval(exchange *Exchange) (any, error) {
	return expr.Run(e.program, exchange.asMap())
}

func (e *simpleExpr) String() string {
	return e.rawExpr
}

type Predicate interface {
	Test(exchange *Exchange) (bool, error)
}

type PredicateFunc func(exchange *Exchange) (bool, error)

func (prd PredicateFunc) Test(exchange *Exchange) (bool, error) {
	return prd(exchange)
}

func newPredicateFromExpr(expr expression) PredicateFunc {
	return func(exchange *Exchange) (bool, error) {
		v, err := expr.eval(exchange)
		if err != nil {
			return false, err
		}
		return toBool(v)
	}
}

func toBool(v any) (bool, error) {
	switch x := v.(type) {
	case nil:
		return false, nil
	case bool:
		return x, nil
	case int:
		return x != 0, nil
	case int8:
		return x != 0, nil
	case int16:
		return x != 0, nil
	case int32:
		return x != 0, nil
	case int64:
		return x != 0, nil
	case uint:
		return x != 0, nil
	case uint8:
		return x != 0, nil
	case uint16:
		return x != 0, nil
	case uint32:
		return x != 0, nil
	case uint64:
		return x != 0, nil
	case float32:
		return x != 0, nil
	case float64:
		return x != 0, nil
	case string:
		if b, err := strconv.ParseBool(x); err == nil {
			return b, nil
		}
		if x == "1" {
			return true, nil
		}
		if x == "0" {
			return false, nil
		}
		return false, fmt.Errorf("cannot convert string %q to bool", x)
	default:
		return false, fmt.Errorf("unsupported type %T", v)
	}
}
