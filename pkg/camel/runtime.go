package camel

import (
	"context"
	"errors"
	"fmt"
	"github.com/paveldanilin/go-camel/pkg/camel/api"
	"github.com/paveldanilin/go-camel/pkg/camel/converter"
	"github.com/paveldanilin/go-camel/pkg/camel/dataformat"
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
	"github.com/paveldanilin/go-camel/pkg/camel/logger"
	"github.com/paveldanilin/go-camel/pkg/camel/uri"
	"log/slog"
	"os"
	"reflect"
	"sync"
)

type EndpointRegistry interface {
	Endpoint(string) api.Endpoint
}

type ExchangeFactoryAware interface {
	SetExchangeFactory(f api.ExchangeFactory)
}

type ConverterRegistry interface {
	// Register registers new converter that MUST implement Converter interface.
	Register(conv any) error
	Type(name string) (reflect.Type, bool)
	Convert(value any, toType reflect.Type, params map[string]any) (any, error)
}

type route struct {
	name     string
	from     string
	producer api.Producer
}

type Runtime struct {
	name string
	mu   sync.RWMutex

	messageHistory bool

	funcRegistry       FuncRegistry
	componentRegistry  ComponentRegistry
	dataFormatRegistry DataFormatRegistry
	exchangeFactory    api.ExchangeFactory
	converterRegistry  ConverterRegistry

	routes    map[string]*route
	endpoints map[string]api.Endpoint
	consumers []api.Consumer

	logger api.Logger

	ctx    context.Context
	cancel context.CancelFunc
}

type RuntimeConfig struct {
	Name               string
	ExchangeFactory    api.ExchangeFactory
	FuncRegistry       FuncRegistry
	ComponentRegistry  ComponentRegistry
	DataFormatRegistry DataFormatRegistry
	ConverterRegistry  ConverterRegistry
	Logger             api.Logger
	MessageHistory     bool
}

func NewRuntime(config RuntimeConfig) *Runtime {
	ctx, cancel := context.WithCancel(context.Background())

	runtime := &Runtime{
		name:               config.Name,
		funcRegistry:       config.FuncRegistry,
		componentRegistry:  config.ComponentRegistry,
		dataFormatRegistry: config.DataFormatRegistry,
		exchangeFactory:    config.ExchangeFactory,
		converterRegistry:  config.ConverterRegistry,
		logger:             config.Logger,

		messageHistory: config.MessageHistory,

		routes:    map[string]*route{},
		endpoints: map[string]api.Endpoint{},
		consumers: []api.Consumer{},

		ctx:    ctx,
		cancel: cancel,
	}

	if runtime.name == "" {
		runtime.name = "CamelRuntime"
	}
	//runtime.ctx = context.WithValue(runtime.ctx, "CamelRuntimeName", runtime.headerName)

	if runtime.funcRegistry == nil {
		runtime.funcRegistry = newFuncRegistry()
	}
	if runtime.componentRegistry == nil {
		runtime.componentRegistry = newComponentRegistry()
	}

	// register default DataFormat registry
	if runtime.dataFormatRegistry == nil {
		runtime.dataFormatRegistry = newDataFormatRegistry()
		runtime.dataFormatRegistry.RegisterDataFormat("json", &dataformat.JSON{})
		runtime.dataFormatRegistry.RegisterDataFormat("xml", &dataformat.XML{})
	}
	if runtime.logger == nil {
		runtime.logger = logger.NewSlog(slog.New(slog.NewTextHandler(os.Stdout, nil)), api.LogLevelInfo)
	}
	if runtime.converterRegistry == nil {
		runtime.converterRegistry = NewConverterRegistry()
		runtime.converterRegistry.Register(converter.StringToBool())
		runtime.converterRegistry.Register(converter.StringToFloat64())
		runtime.converterRegistry.Register(converter.StringToFloat())
		runtime.converterRegistry.Register(converter.StringToInt64())
		runtime.converterRegistry.Register(converter.StringToInt())
		runtime.converterRegistry.Register(converter.StringToDateTime())
	}

	return runtime
}

func (rt *Runtime) Name() string {
	return rt.name
}

// RegisterFunc registers a named fn in the current Runtime.
func (rt *Runtime) RegisterFunc(name string, fn func(*exchange.Exchange)) error {
	return rt.funcRegistry.RegisterFunc(name, fn)
}

func (rt *Runtime) MustRegisterFunc(name string, fn func(e *exchange.Exchange)) {
	err := rt.RegisterFunc(name, fn)
	if err != nil {
		panic(fmt.Errorf("camel: %w", err))
	}
}

// RegisterComponent register the given Component in the current Runtime.

func (rt *Runtime) RegisterComponent(c api.Component) error {
	err := rt.componentRegistry.RegisterComponent(c)
	if err != nil {
		return err
	}
	if ef, isExchangeFactoryAware := c.(ExchangeFactoryAware); isExchangeFactoryAware {
		ef.SetExchangeFactory(rt)
	}
	return nil
}

func (rt *Runtime) MustRegisterComponent(c api.Component) {
	err := rt.RegisterComponent(c)
	if err != nil {
		panic(fmt.Errorf("camel: %w", err))
	}
}

func (rt *Runtime) Component(id string) api.Component {
	return rt.componentRegistry.Component(id)
}

func (rt *Runtime) RegisterDataFormat(name string, dataFormat api.DataFormat) error {
	return rt.dataFormatRegistry.RegisterDataFormat(name, dataFormat)
}

func (rt *Runtime) MustRegisterDataFormat(name string, dataFormat api.DataFormat) {
	err := rt.RegisterDataFormat(name, dataFormat)
	if err != nil {
		panic(fmt.Errorf("camel: %w", err))
	}
}

func (rt *Runtime) RegisterRoute(routeDefinition *Route) error {
	if _, exists := rt.routes[routeDefinition.Name]; exists {
		rt.logger.Error(context.Background(), fmt.Sprintf("Route with headerName '%s' already registered", routeDefinition.Name))
		return errors.New("step already registered: " + routeDefinition.Name)
	}

	r, err := compileRoute(compilerConfig{
		funcRegistry:       rt.funcRegistry,
		logger:             rt.logger,
		dataFormatRegistry: rt.dataFormatRegistry,
		converterRegistry:  rt.converterRegistry,
		endpointRegistry:   rt,
		preProcessor:       rt.preProcessor,
		postProcessor:      rt.postProcessor,
	}, routeDefinition)
	if err != nil {
		rt.logger.Error(context.Background(), "Route compilation failed", slog.String("error", err.Error()))
		return err
	}

	rt.routes[routeDefinition.Name] = r
	rt.logger.Info(context.Background(), fmt.Sprintf("Route '%s' registered and consuming from: '%s'", routeDefinition.Name, routeDefinition.From))

	return nil
}

func (rt *Runtime) MustRegisterRoute(routeDefinition *Route) {
	err := rt.RegisterRoute(routeDefinition)
	if err != nil {
		panic(fmt.Errorf("camel: failed to register step: %w", err))
	}
}

func (rt *Runtime) Endpoint(uri string) api.Endpoint {
	if endpoint, exists := rt.endpoints[uri]; exists {
		return endpoint
	}
	return nil
}

func (rt *Runtime) NewExchange(c context.Context) *exchange.Exchange {
	var newExchange *exchange.Exchange
	if rt.exchangeFactory == nil {
		// Default exchange factory
		newExchange = exchange.NewExchange(c)
	} else {
		newExchange = rt.exchangeFactory.NewExchange(c)
	}
	if rt.messageHistory {
		newExchange.Message().SetHeader(exchange.CamelHeaderMessageHistory, exchange.NewMessageHistory())
	}
	return newExchange
}

func (rt *Runtime) preProcessor(e *exchange.Exchange) {
	//slog.Info("[pre]", "exchange", slog.AnyValue(exchange))
	// TODO: add event firing
}

func (rt *Runtime) postProcessor(e *exchange.Exchange) {
	//slog.Info("[post]", "exchange", slog.AnyValue(exchange))
	// TODO: add event firing

}

func (rt *Runtime) Send(ctx context.Context, uri string, body any, headers map[string]any) (*exchange.Exchange, error) {
	endpoint := rt.Endpoint(uri)
	if endpoint == nil {
		return nil, errors.New("endpoint not found for uri: " + uri)
	}

	producer, err := endpoint.CreateProducer()
	if err != nil {
		return nil, err
	}

	exchangeCopy := rt.NewExchange(ctx)
	exchangeCopy.Message().Body = body
	exchangeCopy.Message().Headers().SetAll(headers)

	producer.Process(exchangeCopy)

	return exchangeCopy, exchangeCopy.Error()
}

func (rt *Runtime) SendBody(ctx context.Context, uri string, body any) (*exchange.Message, error) {
	exchangeCopy, err := rt.Send(ctx, uri, body, nil)
	if err != nil {
		return nil, err
	}
	if exchangeCopy.Error() != nil {
		return nil, exchangeCopy.Error()
	}
	return exchangeCopy.Message(), nil
}

func (rt *Runtime) SendHeaders(ctx context.Context, uri string, headers map[string]any) (*exchange.Message, error) {
	exchangeCopy, err := rt.Send(ctx, uri, nil, headers)
	if err != nil {
		return nil, err
	}
	if exchangeCopy.Error() != nil {
		return nil, exchangeCopy.Error()
	}
	return exchangeCopy.Message(), nil
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

	for _, r := range rt.routes {
		fromUri, err := uri.ParseURI(r.from, nil)
		if err != nil {
			return fmt.Errorf("invalid URI format in step '%s' that consumes from '%s': %w", r.name, r.from, err)
		}

		// Resolve component
		component := rt.componentRegistry.Component(fromUri.Component())
		if component == nil {
			return fmt.Errorf("component '%s' not found in step '%s' that consumes from '%s'", fromUri.Component(), r.name, r.from)
		}

		// Resolve/create endpoint
		endpoint, endpointExists := rt.endpoints[r.from]
		if !endpointExists {
			endpoint, err = component.CreateEndpoint(r.from)
			if err != nil {
				return fmt.Errorf("failed to create endpoint in step '%s' that consumes from '%s': %w", r.name, r.from, err)
			}
			rt.endpoints[r.from] = endpoint
		}

		// Create consumer
		consumer, err := endpoint.CreateConsumer(r.producer)
		if err != nil {
			return fmt.Errorf("failed to create consumer in step '%s' that consumes from '%s': %w", r.name, r.from, err)
		}
		rt.consumers = append(rt.consumers, consumer)
	}

	// Start consumers
	for _, consumer := range rt.consumers {
		if err := consumer.Start(); err != nil {
			return err
		}
	}

	rt.logger.Info(context.Background(), fmt.Sprintf("Camel runtime '%s' started", rt.name))

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
	rt.routes = nil

	rt.logger.Info(context.Background(), fmt.Sprintf("Camel runetime '%s' stopped", rt.name))

	return nil
}
