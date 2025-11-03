package camel

import (
	"errors"
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
	"sync"
)

type FuncRegistry interface {
	RegisterFunc(name string, fn func(*exchange.Exchange)) error
	Func(name string) func(*exchange.Exchange)
}

type funcRegistry struct {
	mu      sync.Mutex
	funcMap map[string]func(*exchange.Exchange)
}

func newFuncRegistry() *funcRegistry {
	return &funcRegistry{
		funcMap: map[string]func(*exchange.Exchange){},
	}
}

func (fr *funcRegistry) RegisterFunc(name string, fn func(*exchange.Exchange)) error {
	fr.mu.Lock()
	defer fr.mu.Unlock()

	if _, exists := fr.funcMap[name]; exists {
		return errors.New("fn already registered")
	}

	fr.funcMap[name] = fn
	return nil
}

func (fr *funcRegistry) Func(name string) func(*exchange.Exchange) {
	fr.mu.Lock()
	defer fr.mu.Unlock()

	if fn, exists := fr.funcMap[name]; exists {
		return fn
	}

	return nil
}
