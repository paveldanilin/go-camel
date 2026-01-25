package camel

import (
	"errors"
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
	"sync"
)

type Func func(*exchange.Exchange)

type funcRegistry struct {
	mu      sync.Mutex
	funcMap map[string]Func
}

func newFuncRegistry() *funcRegistry {
	return &funcRegistry{
		funcMap: map[string]Func{},
	}
}

func (fr *funcRegistry) RegisterFunc(name string, fn Func) error {
	fr.mu.Lock()
	defer fr.mu.Unlock()

	if _, exists := fr.funcMap[name]; exists {
		return errors.New("fn already registered")
	}

	fr.funcMap[name] = fn
	return nil
}

func (fr *funcRegistry) Func(name string) Func {
	fr.mu.Lock()
	defer fr.mu.Unlock()

	if fn, exists := fr.funcMap[name]; exists {
		return fn
	}

	return nil
}
