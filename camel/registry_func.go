package camel

import (
	"errors"
	"sync"
)

type FuncRegistry interface {
	RegisterFunc(name string, fn func(*Exchange)) error
	Func(name string) func(*Exchange)
}

type funcRegistry struct {
	mu      sync.Mutex
	funcMap map[string]func(*Exchange)
}

func newFuncRegistry() *funcRegistry {
	return &funcRegistry{
		funcMap: map[string]func(*Exchange){},
	}
}

func (fr *funcRegistry) RegisterFunc(name string, fn func(*Exchange)) error {
	fr.mu.Lock()
	defer fr.mu.Unlock()

	if _, exists := fr.funcMap[name]; exists {
		return errors.New("func already registered")
	}

	fr.funcMap[name] = fn
	return nil
}

func (fr *funcRegistry) Func(name string) func(*Exchange) {
	fr.mu.Lock()
	defer fr.mu.Unlock()

	if fn, exists := fr.funcMap[name]; exists {
		return fn
	}

	return nil
}
