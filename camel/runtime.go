package camel

import (
	"context"
	"errors"
	"fmt"
	"github.com/paveldanilin/go-camel/camel/component"
	"sync"
)

type Processor interface {
	Process(message *Message)
}

type Consumer interface {
	Start() error
	Stop() error
}

type Producer interface {
	Processor
}

type Endpoint interface {
	Uri() string
	CreateConsumer(processor Processor) (Consumer, error)
	CreateProducer() (Producer, error)
}

type Component interface {
	Id() string
	CreateEndpoint(uri string) (Endpoint, error)
}

type RuntimeAware interface {
	SetRuntime(runtime *Runtime)
}

type Runtime struct {
	components map[string]Component
	routes     map[string]*Route
	endpoints  map[string]Endpoint
	consumers  []Consumer
	mu         sync.RWMutex
	ctx        context.Context
	cancel     context.CancelFunc
}

func NewRuntime() *Runtime {
	ctx, cancel := context.WithCancel(context.Background())
	return &Runtime{
		components: map[string]Component{},
		routes:     map[string]*Route{},
		endpoints:  map[string]Endpoint{},
		consumers:  []Consumer{},
		ctx:        ctx,
		cancel:     cancel,
	}
}

func (ctx *Runtime) RegisterComponent(c Component) error {

	if _, exists := ctx.components[c.Id()]; exists {
		return errors.New("component already registered: " + c.Id())
	}

	if contextAware, isContextAware := c.(RuntimeAware); isContextAware {
		contextAware.SetRuntime(ctx)
	}

	ctx.components[c.Id()] = c

	return nil
}

func (ctx *Runtime) Component(componentId string) Component {

	if c, exists := ctx.components[componentId]; exists {
		return c
	}

	return nil
}

func (ctx *Runtime) RegisterRoute(r *Route) error {

	if _, exists := ctx.routes[r.id]; exists {
		return errors.New("route already registered: " + r.id)
	}

	ctx.routes[r.id] = r

	return nil
}

func (ctx *Runtime) Endpoint(uri string) Endpoint {

	return ctx.endpoints[uri]
}

func (ctx *Runtime) Send(uri string, payload any, headers map[string]any) (*Message, error) {

	if endpoint, exists := ctx.endpoints[uri]; exists {
		producer, err := endpoint.CreateProducer()
		if err != nil {
			return nil, err
		}

		// TODO: message pooling?
		message := NewMessage()
		message.runtime = ctx
		message.payload = payload
		message.headers.SetAll(headers)

		producer.Process(message)

		return message, nil
	}

	return nil, errors.New("endpoint not found for uri: " + uri)
}

func (ctx *Runtime) Route(routeId string) *Route {

	if r, exists := ctx.routes[routeId]; exists {
		return r
	}

	return nil
}

func (ctx *Runtime) Start() error {

	ctx.mu.Lock()
	defer ctx.mu.Unlock()

	for _, route := range ctx.routes {

		parts := component.SplitURI(route.from)
		if len(parts) != 2 {
			return fmt.Errorf("invalid URI format: %s", route.from)
		}

		componentName, uri := parts[0], parts[1]
		component, exists := ctx.components[componentName]
		if !exists {
			return fmt.Errorf("component not found: %s", componentName)
		}

		if endpoint, exists := ctx.endpoints[route.from]; exists {
			_, err := endpoint.CreateConsumer(route.producer)
			if err != nil {
				return fmt.Errorf("zzz: %w", err)
			}
		} else {
			endpoint, err := component.CreateEndpoint(uri)
			if err != nil {
				return fmt.Errorf("failed to create endpoint: %w", err)
			}
			ctx.endpoints[route.from] = endpoint
			consumer, err := endpoint.CreateConsumer(route.producer)
			if err != nil {
				return fmt.Errorf("zzz: %w", err)
			}
			ctx.consumers = append(ctx.consumers, consumer)
		}
	}

	for _, consumer := range ctx.consumers {
		if err := consumer.Start(); err != nil {
			return err
		}
	}

	return nil

}

func (ctx *Runtime) Stop() error {

	ctx.mu.Lock()
	defer ctx.mu.Unlock()

	ctx.cancel()

	for _, consumer := range ctx.consumers {
		if err := consumer.Stop(); err != nil {
			return err
		}
	}

	ctx.consumers = nil
	ctx.endpoints = nil
	ctx.components = nil
	ctx.routes = nil

	return nil
}
