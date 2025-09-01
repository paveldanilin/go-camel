package camel

import (
	"context"
	"errors"
	"fmt"
	"github.com/paveldanilin/go-camel/camel/uri"
	"log"
	"sync"
)

type Expr interface {
	Eval(exchange *Exchange) (any, error)
}

type Processor interface {
	Process(exchange *Exchange)
}

type Consumer interface {
	Start() error
	Stop() error
}

type Producer interface {
	Processor
}

type Endpoint interface {
	Uri() *uri.URI
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

func (rt *Runtime) RegisterComponent(c Component) error {
	if _, exists := rt.components[c.Id()]; exists {
		return errors.New("component already registered: " + c.Id())
	}

	if contextAware, isContextAware := c.(RuntimeAware); isContextAware {
		contextAware.SetRuntime(rt)
	}

	rt.components[c.Id()] = c

	return nil
}

func (rt *Runtime) MustRegisterComponent(c Component) {
	err := rt.RegisterComponent(c)
	if err != nil {
		panic(err)
	}
}

func (rt *Runtime) Component(componentId string) Component {
	if c, exists := rt.components[componentId]; exists {
		return c
	}

	return nil
}

func (rt *Runtime) RegisterRoute(r *Route) error {
	if _, exists := rt.routes[r.id]; exists {
		return errors.New("route already registered: " + r.id)
	}

	rt.routes[r.id] = r

	return nil
}

func (rt *Runtime) MustRegisterRoute(r *Route) {
	err := rt.RegisterRoute(r)
	if err != nil {
		panic(err)
	}
}

func (rt *Runtime) Endpoint(uri string) Endpoint {
	if endpoint, exists := rt.endpoints[uri]; exists {
		return endpoint
	}
	return nil
}

func (rt *Runtime) Send(uri string, body any, headers map[string]any) (*Message, error) {
	if endpoint, exists := rt.endpoints[uri]; exists {
		producer, err := endpoint.CreateProducer()
		if err != nil {
			return nil, err
		}

		e := NewExchange(nil, rt)
		e.message.Body = body
		e.message.headers.SetAll(headers)

		producer.Process(e)

		return e.message, nil
	}

	return nil, errors.New("endpoint not found for uri: " + uri)
}

func (rt *Runtime) Route(routeId string) *Route {
	if r, exists := rt.routes[routeId]; exists {
		return r
	}

	return nil
}

func (rt *Runtime) Start() error {
	rt.mu.Lock()
	defer rt.mu.Unlock()

	for _, route := range rt.routes {
		log.Printf("Registering route: '%s'", route.id)

		fromUri, err := uri.Parse(route.from, nil)
		if err != nil {
			return fmt.Errorf("invalid URI format in route '%s' that consumes from [%s]: %w", route.id, route.from, err)
		}

		log.Printf("Route uri parsed: %s", fromUri)

		// Resolve component
		component, componentExists := rt.components[fromUri.Component()]
		if !componentExists {
			return fmt.Errorf("component '%s' not found in route '%s' that consumes from [%s]", fromUri.Component(), route.id, route.from)
		}

		log.Printf("Compoentn resolved: %s", fromUri.Component())

		// Resolve/create endpoint
		endpoint, endpointExists := rt.endpoints[route.from]
		if !endpointExists {
			endpoint, err = component.CreateEndpoint(route.from)
			if err != nil {
				return fmt.Errorf("failed to create endpoint in route '%s' that consumes from [%s]: %w", route.id, route.from, err)
			}
			rt.endpoints[route.from] = endpoint
		}

		log.Printf("Endpoint created")

		// Create consumer
		consumer, err := endpoint.CreateConsumer(route.producer)
		if err != nil {
			return fmt.Errorf("failed to create consumer in route '%s' that consumes from [%s]: %w", route.id, route.from, err)
		}
		rt.consumers = append(rt.consumers, consumer)
	}

	// Start consumers
	for _, consumer := range rt.consumers {
		if err := consumer.Start(); err != nil {
			return err
		}
	}

	return nil
}

func (rt *Runtime) Stop() error {
	rt.mu.Lock()
	defer rt.mu.Unlock()

	rt.cancel()

	for _, consumer := range rt.consumers {
		if err := consumer.Stop(); err != nil {
			return err
		}
	}

	rt.consumers = nil
	rt.endpoints = nil
	rt.components = nil
	rt.routes = nil

	return nil
}
