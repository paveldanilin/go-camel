package camel

import (
	"errors"
	"github.com/paveldanilin/go-camel/pkg/camel/api"
	"sync"
)

type ComponentRegistry interface {
	RegisterComponent(component api.Component) error
	Component(id string) api.Component
}

type componentRegistry struct {
	mu           sync.Mutex
	componentMap map[string]api.Component
}

func newComponentRegistry() *componentRegistry {
	return &componentRegistry{
		componentMap: map[string]api.Component{},
	}
}

func (cr *componentRegistry) RegisterComponent(component api.Component) error {
	cr.mu.Lock()
	defer cr.mu.Unlock()

	if _, exists := cr.componentMap[component.Id()]; exists {
		return errors.New("component already registered")
	}

	cr.componentMap[component.Id()] = component
	return nil
}

func (cr *componentRegistry) Component(id string) api.Component {
	cr.mu.Lock()
	defer cr.mu.Unlock()

	if c, exists := cr.componentMap[id]; exists {
		return c
	}

	return nil
}
