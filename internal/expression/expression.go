package expression

import (
	"fmt"
	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
	"strconv"
)

// Expression represents an expression that takes Exchange and returns computed valueExpression or error.
// Used to dynamically compute valueExpression to setBodyProcessor/setHeaderProcessor.
type Expression interface {
	Eval(e *exchange.Exchange) (any, error)
}

// Func represents a user function that returns valueExpression.
type Func func(e *exchange.Exchange) (any, error)

func (fn Func) Eval(e *exchange.Exchange) (any, error) {
	return fn(e)
}

// constant represents a constant value.
type constant struct {
	value any
}

func NewConst(value any) *constant {
	return &constant{
		value: value,
	}
}

func (e *constant) Eval(_ *exchange.Exchange) (any, error) {
	return e.value, nil
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
		panic(fmt.Errorf("camel: expr: simple: %w", err))
	}

	return expr_
}

func (e *simple) Eval(ex *exchange.Exchange) (any, error) {
	return expr.Run(e.program, ex.AsMap())
}

func (e *simple) String() string {
	return e.raw
}

type Predicate interface {
	Test(e *exchange.Exchange) (bool, error)
}

type PredicateFunc func(e *exchange.Exchange) (bool, error)

func (prd PredicateFunc) Test(e *exchange.Exchange) (bool, error) {
	return prd(e)
}

func NewPredicateFromExpression(expr Expression) PredicateFunc {
	return func(e *exchange.Exchange) (bool, error) {
		v, err := expr.Eval(e)
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
