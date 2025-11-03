package api

import (
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
	"github.com/paveldanilin/go-camel/pkg/camel/uri"
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
	Id() string
	CreateEndpoint(uri string) (Endpoint, error)
}

type DataFormat interface {
	Unmarshal(data []byte, targetType any) (any, error)
	Marshal(data any) (string, error)
}

type Converter[From any, To any] interface {
	Convert(from From, params map[string]any) (To, error)
}

type RouteStep interface {
	StepName() string
}
