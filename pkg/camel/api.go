package camel

import (
	"context"
	"github.com/paveldanilin/go-camel/internal/processor"
	"github.com/paveldanilin/go-camel/pkg/camel/dataformat"
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
	"github.com/paveldanilin/go-camel/pkg/camel/uri"
	"reflect"
)

type Processor interface {
	Process(e *exchange.Exchange)
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
	ComponentId() string
	CreateEndpoint(uri string) (Endpoint, error)
}

type EndpointRegistry interface {
	CreateEndpointProducer(uri string) (processor.Processor, error)
}

type ExchangeFactory interface {
	NewExchange(c context.Context) *exchange.Exchange
}

type ExchangeFactoryAware interface {
	SetExchangeFactory(f ExchangeFactory)
}

type Converter[From any, To any] interface {
	Convert(from From, params map[string]any) (To, error)
}

type ConverterRegistry interface {
	// Register registers new converter that MUST implement Converter interface.
	Register(conv any) error
	Type(name string) (reflect.Type, bool)
	Convert(value any, toType reflect.Type, params map[string]any) (any, error)
}

type ComponentRegistry interface {
	RegisterComponent(component Component) error
	Component(id string) Component
}

type DataFormatRegistry interface {
	RegisterDataFormat(name string, format dataformat.DataFormat) error
	DataFormat(name string) dataformat.DataFormat
}

type FuncRegistry interface {
	RegisterFunc(name string, fn Func) error
	Func(name string) Func
}

type Env interface {
	LookupVar(name string) (string, bool)
}
