package expr

import "github.com/paveldanilin/go-camel/pkg/camel/exchange"

type Kind string

const (
	SimpleKind   Kind = "simple"
	ConstantKind      = "constant"
	FuncKind          = "func"
)

type Definition struct {
	Kind       Kind
	Expression any
}

func Simple(expression string) Definition {
	return Definition{
		Kind:       SimpleKind,
		Expression: expression,
	}
}

func Constant(value any) Definition {
	return Definition{
		Kind:       ConstantKind,
		Expression: value,
	}
}

func Func(fn func(e *exchange.Exchange) (any, error)) Definition {
	return Definition{
		Kind:       FuncKind,
		Expression: fn,
	}
}
