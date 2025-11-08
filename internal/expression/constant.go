package expression

import "github.com/paveldanilin/go-camel/pkg/camel/exchange"

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
