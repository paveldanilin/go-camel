package delay

import (
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
	"time"
)

type delayProcessor struct {
	routeName string
	name      string
	durMs     int64
}

func NewProcessor(routeName, name string, durMs int64) *delayProcessor {
	return &delayProcessor{
		routeName: routeName,
		name:      name,
		durMs:     durMs,
	}
}

func (p *delayProcessor) Name() string {
	return p.name
}

func (p *delayProcessor) RouteName() string {
	return p.routeName
}

func (p *delayProcessor) Process(_ *exchange.Exchange) {
	time.Sleep(time.Duration(p.durMs) * time.Millisecond)
}
