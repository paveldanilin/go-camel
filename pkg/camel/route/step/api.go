package step

import "github.com/paveldanilin/go-camel/pkg/camel/exchange"

type RouteStep interface {
	StepName() string
}

type ExchangeAggregator interface {
	AggregateExchange(oldExchange, newExchange *exchange.Exchange) *exchange.Exchange
}
