package processor

import (
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
	"reflect"
)

type Processor interface {
	Process(e *exchange.Exchange)
}

type P func(e *exchange.Exchange)

func (f P) Process(e *exchange.Exchange) {
	f(e)
}

type Converter interface {
	Convert(any, reflect.Type, map[string]any) (any, error)
}

type DataFormat interface {
	Unmarshal(data []byte, targetType any) (any, error)
	Marshal(data any) (string, error)
}

type Expression interface {
	Eval(e *exchange.Exchange) (any, error)
}

type ExpressionFunc func(e *exchange.Exchange) (any, error)

func (f ExpressionFunc) Eval(e *exchange.Exchange) (any, error) {
	return f(e)
}

func constant(v any) ExpressionFunc {
	return func(e *exchange.Exchange) (any, error) {
		return v, nil
	}
}

type Predicate interface {
	Test(e *exchange.Exchange) (bool, error)
}

type PredicateFunc func(e *exchange.Exchange) (bool, error)

func (f PredicateFunc) Test(e *exchange.Exchange) (bool, error) {
	return f(e)
}
