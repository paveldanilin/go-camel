package expression

import "github.com/paveldanilin/go-camel/pkg/camel/exchange"

type Constant struct {
	value any
}

func NewConstant(value any) *Constant {
	return &Constant{
		value: value,
	}
}

func (e *Constant) Eval(_ *exchange.Exchange) (any, error) {
	return e.value, nil
}
