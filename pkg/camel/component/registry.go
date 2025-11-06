package component

import (
	"errors"
	"github.com/paveldanilin/go-camel/pkg/camel/api"
	"sync"
)

type registry struct {
	mu           sync.Mutex
	componentMap map[string]api.Component
}

func NewRegistry() *registry {
	return &registry{
		componentMap: map[string]api.Component{},
	}
}

func (cr *registry) RegisterComponent(component api.Component) error {
	cr.mu.Lock()
	defer cr.mu.Unlock()

	if _, exists := cr.componentMap[component.Id()]; exists {
		return errors.New("component already registered")
	}

	cr.componentMap[component.Id()] = component
	return nil
}

func (cr *registry) Component(id string) api.Component {
	cr.mu.Lock()
	defer cr.mu.Unlock()

	if c, exists := cr.componentMap[id]; exists {
		return c
	}

	return nil
}
