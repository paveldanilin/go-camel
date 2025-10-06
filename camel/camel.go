package camel

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"sync"
)

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
	Uri() *URI
	CreateConsumer(processor Processor) (Consumer, error)
	CreateProducer() (Producer, error)
}

type Component interface {
	Id() string
	CreateEndpoint(uri string) (Endpoint, error)
}

// RuntimeAware indicates that Component implementation uses Runtime internally.
type RuntimeAware interface {
	SetRuntime(runtime *Runtime)
}

type route struct {
	name     string
	from     string
	producer Producer
}

type Runtime struct {
	mu sync.RWMutex

	funcRegistry       FuncRegistry
	componentRegistry  ComponentRegistry
	dataFormatRegistry DataFormatRegistry
	exchangeFactory    ExchangeFactory

	routes    map[string]*route
	endpoints map[string]Endpoint
	consumers []Consumer

	logger Logger

	ctx    context.Context
	cancel context.CancelFunc
}

type RuntimeConfig struct {
	ExchangeFactory    ExchangeFactory
	FuncRegistry       FuncRegistry
	ComponentRegistry  ComponentRegistry
	DataFormatRegistry DataFormatRegistry
	Logger             Logger
}

func NewRuntime(config RuntimeConfig) *Runtime {
	ctx, cancel := context.WithCancel(context.Background())

	runtime := &Runtime{
		funcRegistry:       config.FuncRegistry,
		componentRegistry:  config.ComponentRegistry,
		dataFormatRegistry: config.DataFormatRegistry,
		exchangeFactory:    config.ExchangeFactory,
		logger:             config.Logger,

		routes:    map[string]*route{},
		endpoints: map[string]Endpoint{},
		consumers: []Consumer{},

		ctx:    ctx,
		cancel: cancel,
	}

	if runtime.funcRegistry == nil {
		runtime.funcRegistry = newFuncRegistry()
	}
	if runtime.componentRegistry == nil {
		runtime.componentRegistry = newComponentRegistry()
	}
	if runtime.dataFormatRegistry == nil {
		runtime.dataFormatRegistry = newDataFormatRegistry()
	}
	if runtime.logger == nil {
		runtime.logger = NewSlogLogger(slog.New(slog.NewTextHandler(os.Stdout, nil)), LogLevelInfo)
	}

	return runtime
}

// RegisterFunc registers a named func in the current Runtime.
func (rt *Runtime) RegisterFunc(name string, fn func(*Exchange)) error {
	return rt.funcRegistry.RegisterFunc(name, fn)
}

func (rt *Runtime) MustRegisterFunc(name string, fn func(exchange *Exchange)) {
	err := rt.RegisterFunc(name, fn)
	if err != nil {
		panic(fmt.Errorf("camel: %w", err))
	}
}

// RegisterComponent register the given Component in the current Runtime.
func (rt *Runtime) RegisterComponent(c Component) error {
	err := rt.componentRegistry.RegisterComponent(c)
	if err != nil {
		return err
	}
	if rtAware, isRtAware := c.(RuntimeAware); isRtAware {
		rtAware.SetRuntime(rt)
	}
	return nil
}

func (rt *Runtime) MustRegisterComponent(c Component) {
	err := rt.RegisterComponent(c)
	if err != nil {
		panic(fmt.Errorf("camel: %w", err))
	}
}

func (rt *Runtime) Component(id string) Component {
	return rt.componentRegistry.Component(id)
}

func (rt *Runtime) RegisterDataFormat(name string, dataFormat DataFormat) error {
	return rt.dataFormatRegistry.RegisterDataFormat(name, dataFormat)
}

func (rt *Runtime) MustRegisterDataFormat(name string, dataFormat DataFormat) {
	err := rt.RegisterDataFormat(name, dataFormat)
	if err != nil {
		panic(fmt.Errorf("camel: %w", err))
	}
}

func (rt *Runtime) RegisterRoute(routeDefinition *Route) error {
	if _, exists := rt.routes[routeDefinition.Name]; exists {
		return errors.New("route already registered: " + routeDefinition.Name)
	}

	r, err := compileRoute(compilerConfig{
		funcRegistry:       rt.funcRegistry,
		preProcessorFunc:   rt.preProcessor,
		postProcessorFunc:  rt.postProcessor,
		logger:             rt.logger,
		dataFormatRegistry: rt.dataFormatRegistry,
	}, routeDefinition)
	if err != nil {
		return err
	}

	rt.routes[routeDefinition.Name] = r

	return nil
}

func (rt *Runtime) MustRegisterRoute(routeDefinition *Route) {
	err := rt.RegisterRoute(routeDefinition)
	if err != nil {
		panic(fmt.Errorf("camel: failed to register route: %w", err))
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
	//slog.Info("[pre]", "exchange", slog.AnyValue(exchange))
	// TODO: add event firing
}

func (rt *Runtime) postProcessor(exchange *Exchange) {
	//slog.Info("[post]", "exchange", slog.AnyValue(exchange))
	// TODO: add event firing

}

func (rt *Runtime) Send(ctx context.Context, uri string, body any, headers map[string]any) (*Exchange, error) {
	endpoint := rt.Endpoint(uri)
	if endpoint == nil {
		return nil, errors.New("endpoint not found for uri: " + uri)
	}

	producer, err := endpoint.CreateProducer()
	if err != nil {
		return nil, err
	}

	exchange := rt.NewExchange(ctx)
	//exchange.SetProperty("CAMEL_ROUTE_NAME", "")
	exchange.message.Body = body
	exchange.message.headers.SetAll(headers)

	producer.Process(exchange)

	return exchange, exchange.Error()
}

func (rt *Runtime) SendBody(ctx context.Context, uri string, body any) (*Message, error) {
	exchange, err := rt.Send(ctx, uri, body, nil)
	if err != nil {
		return nil, err
	}
	if exchange.err != nil {
		return nil, exchange.err
	}
	return exchange.message, nil
}

func (rt *Runtime) SendHeaders(ctx context.Context, uri string, headers map[string]any) (*Message, error) {
	exchange, err := rt.Send(ctx, uri, nil, headers)
	if err != nil {
		return nil, err
	}
	if exchange.err != nil {
		return nil, exchange.err
	}
	return exchange.message, nil
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

		fromUri, err := Parse(route.from, nil)
		if err != nil {
			return fmt.Errorf("invalid URI format in route '%s' that consumes from '%s': %w", route.name, route.from, err)
		}

		// Resolve component
		component := rt.componentRegistry.Component(fromUri.Component())
		if component == nil {
			return fmt.Errorf("component '%s' not found in route '%s' that consumes from '%s'", fromUri.Component(), route.name, route.from)
		}

		// Resolve/create endpoint
		endpoint, endpointExists := rt.endpoints[route.from]
		if !endpointExists {
			endpoint, err = component.CreateEndpoint(route.from)
			if err != nil {
				return fmt.Errorf("failed to create endpoint in route '%s' that consumes from '%s': %w", route.name, route.from, err)
			}
			rt.endpoints[route.from] = endpoint
		}

		// Create consumer
		consumer, err := endpoint.CreateConsumer(route.producer)
		if err != nil {
			return fmt.Errorf("failed to create consumer in route '%s' that consumes from '%s': %w", route.name, route.from, err)
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
	//rt.components = nil
	rt.routes = nil

	return nil
}
