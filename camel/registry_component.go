package camel

import (
	"errors"
	"sync"
)

type ComponentRegistry interface {
	RegisterComponent(component Component) error
	Component(id string) Component
}

type componentRegistry struct {
	mu           sync.Mutex
	componentMap map[string]Component
}

func newComponentRegistry() *componentRegistry {
	return &componentRegistry{
		componentMap: map[string]Component{},
	}
}

func (cr *componentRegistry) RegisterComponent(component Component) error {
	cr.mu.Lock()
	defer cr.mu.Unlock()

	if _, exists := cr.componentMap[component.Id()]; exists {
		return errors.New("component already registered")
	}

	cr.componentMap[component.Id()] = component
	return nil
}

func (cr *componentRegistry) Component(id string) Component {
	cr.mu.Lock()
	defer cr.mu.Unlock()

	if c, exists := cr.componentMap[id]; exists {
		return c
	}

	return nil
}
