package camel

import (
	"context"
	"errors"
	"fmt"
	"github.com/paveldanilin/go-camel/camel/dsl"
	"github.com/paveldanilin/go-camel/camel/uri"
	"log/slog"
	"sync"
)

type Expr interface {
	Eval(exchange *Exchange) (any, error)
}

type Predicate interface {
	Test(exchange *Exchange) (bool, error)
}

type PredicateFunc func(exchange *Exchange) (bool, error)

func (prd PredicateFunc) Test(exchange *Exchange) (bool, error) {
	return prd(exchange)
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

type route struct {
	name     string
	from     string
	producer Producer
}

type Runtime struct {
	mu              sync.RWMutex
	exchangeFactory ExchangeFactory
	components      map[string]Component
	routes          map[string]*route
	endpoints       map[string]Endpoint
	consumers       []Consumer
	ctx             context.Context
	cancel          context.CancelFunc
}

type RuntimeOption func(*Runtime)

func NewRuntime(options ...RuntimeOption) *Runtime {
	ctx, cancel := context.WithCancel(context.Background())

	runtime := &Runtime{
		components: map[string]Component{},
		routes:     map[string]*route{},
		endpoints:  map[string]Endpoint{},
		consumers:  []Consumer{},
		ctx:        ctx,
		cancel:     cancel,
	}

	for _, opt := range options {
		opt(runtime)
	}

	return runtime
}

func RuntimeWithExchangeFactory(exchangeFactory ExchangeFactory) RuntimeOption {
	return func(r *Runtime) {
		r.exchangeFactory = exchangeFactory
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
		panic(fmt.Errorf("camel: %w", err))
	}
}

func (rt *Runtime) Component(id string) Component {
	if c, exists := rt.components[id]; exists {
		return c
	}

	return nil
}

func (rt *Runtime) RegisterRoute(r *dsl.Route) error {
	if _, exists := rt.routes[r.Name]; exists {
		return errors.New("route already registered: " + r.Name)
	}

	rtRoute, err := compileRoute(r, compilerConfig{
		preProcessorFunc:  rt.preProcessor,
		postProcessorFunc: rt.postProcessor,
	})
	if err != nil {
		return err
	}

	rt.routes[r.Name] = rtRoute

	return nil
}

func (rt *Runtime) MustRegisterRoute(r *dsl.Route) {
	err := rt.RegisterRoute(r)
	if err != nil {
		panic(fmt.Errorf("camel: %w", err))
	}
}

func (rt *Runtime) Endpoint(uri string) Endpoint {
	if endpoint, exists := rt.endpoints[uri]; exists {
		return endpoint
	}
	return nil
}

func (rt *Runtime) NewExchange(c context.Context) *Exchange {
	var newExchange *Exchange
	if rt.exchangeFactory == nil {
		newExchange = NewExchange(c, rt)
	} else {
		newExchange = rt.exchangeFactory.NewExchange(c)
	}

	newExchange.runtime = rt

	return newExchange
}

func (rt *Runtime) preProcessor(exchange *Exchange) {
	slog.Info("[pre]", "exchange", slog.AnyValue(exchange))
}

func (rt *Runtime) postProcessor(exchange *Exchange) {
	slog.Info("[post]", "exchange", slog.AnyValue(exchange))
}

func (rt *Runtime) Send(ctx context.Context, uri string, body any, headers map[string]any) (*Message, error) {
	endpoint := rt.Endpoint(uri)
	if endpoint == nil {
		return nil, errors.New("endpoint not found for uri: " + uri)
	}

	producer, err := endpoint.CreateProducer()
	if err != nil {
		return nil, err
	}

	ex := rt.NewExchange(ctx)
	ex.message.Body = body
	ex.message.headers.SetAll(headers)

	producer.Process(ex)

	return ex.message, ex.Error()
}

func (rt *Runtime) SendBody(ctx context.Context, uri string, body any) (*Message, error) {
	return rt.Send(ctx, uri, body, nil)
}

func (rt *Runtime) SendHeaders(ctx context.Context, uri string, headers map[string]any) (*Message, error) {
	return rt.Send(ctx, uri, nil, headers)
}

func (rt *Runtime) Route(routeId string) *route {
	if r, exists := rt.routes[routeId]; exists {
		return r
	}

	return nil
}

func (rt *Runtime) Start() error {
	rt.mu.Lock()
	defer rt.mu.Unlock()

	for _, route := range rt.routes {

		fromUri, err := uri.Parse(route.from, nil)
		if err != nil {
			return fmt.Errorf("invalid URI format in dsl '%s' that consumes from [%s]: %w", route.name, route.from, err)
		}

		// Resolve component
		component, componentExists := rt.components[fromUri.Component()]
		if !componentExists {
			return fmt.Errorf("component '%s' not found in dsl '%s' that consumes from [%s]", fromUri.Component(), route.name, route.from)
		}

		// Resolve/create endpoint
		endpoint, endpointExists := rt.endpoints[route.from]
		if !endpointExists {
			endpoint, err = component.CreateEndpoint(route.from)
			if err != nil {
				return fmt.Errorf("failed to create endpoint in dsl '%s' that consumes from [%s]: %w", route.name, route.from, err)
			}
			rt.endpoints[route.from] = endpoint
		}

		// Create consumer
		consumer, err := endpoint.CreateConsumer(route.producer)
		if err != nil {
			return fmt.Errorf("failed to create consumer in dsl '%s' that consumes from [%s]: %w", route.name, route.from, err)
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
