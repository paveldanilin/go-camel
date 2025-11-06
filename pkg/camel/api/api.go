package api

import (
	"context"
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

type LogLevel int

const (
	LogLevelError = iota + 1
	LogLevelWarn
	LogLevelInfo
	LogLevelDebug
)

type Logger interface {
	Log(ctx context.Context, level LogLevel, msg string, args ...any)
	Info(ctx context.Context, msg string, args ...any)
	Warn(ctx context.Context, msg string, args ...any)
	Error(ctx context.Context, msg string, args ...any)
	Debug(ctx context.Context, msg string, args ...any)
}

type ExchangeAggregator interface {
	AggregateExchange(oldExchange *exchange.Exchange, newExchange *exchange.Exchange) *exchange.Exchange
}

type ExchangeFactory interface {
	NewExchange(c context.Context) *exchange.Exchange
}
