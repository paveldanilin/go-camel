package processor

import (
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
	"time"
)

type Delay struct {
	durMs int64
}

func NewDelay(durMs int64) *Delay {
	return &Delay{
		durMs: durMs,
	}
}

func (p *Delay) Process(_ *exchange.Exchange) {
	time.Sleep(time.Duration(p.durMs) * time.Millisecond)
}
