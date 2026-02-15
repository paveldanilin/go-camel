package camel

import (
	"errors"
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
)

type Func func(*exchange.Exchange)

type funcRegistry struct {
	funcMap map[string]Func
}

func newFuncRegistry() *funcRegistry {
	return &funcRegistry{
		funcMap: map[string]Func{},
	}
}

func (fr *funcRegistry) RegisterFunc(name string, fn Func) error {
	if _, exists := fr.funcMap[name]; exists {
		return errors.New("fn already registered")
	}

	fr.funcMap[name] = fn
	return nil
}

func (fr *funcRegistry) Func(name string) Func {
	if fn, exists := fr.funcMap[name]; exists {
		return fn
	}

	return nil
}
