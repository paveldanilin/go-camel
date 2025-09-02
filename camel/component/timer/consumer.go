package timer

import (
	"github.com/paveldanilin/go-camel/camel"
	"sync"
	"time"
)

const (
	HeaderTimerName     = "CamelTimerName"
	HeaderTimerCounter  = "CamelTimerCounter"
	HeaderTimeFiredTime = "CamelTimerFiredTime"
)

type Consumer struct {
	mu         sync.Mutex
	done       chan struct{}
	running    bool
	endpoint   *Endpoint
	processors []camel.Processor
	ticker     *time.Ticker
}

func NewConsumer(endpoint *Endpoint) (*Consumer, error) {
	return &Consumer{
		endpoint:   endpoint,
		processors: []camel.Processor{},
	}, nil
}

func (c *Consumer) Start() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.running {
		return nil
	}

	c.ticker = time.NewTicker(c.endpoint.interval)
	c.done = make(chan struct{})
	c.running = true

	go func() {
		count := int64(0)

		for {
			select {
			case <-c.done:
				return
			case t := <-c.ticker.C:
				count++
				for _, processor := range c.processors {
					exchange := c.endpoint.component.runtime.NewExchange(nil)
					exchange.Message().SetHeader(HeaderTimeFiredTime, t)
					exchange.Message().SetHeader(HeaderTimerName, c.endpoint.name)
					exchange.Message().SetHeader(HeaderTimerCounter, count)

					processor.Process(exchange)
				}
			}
		}
	}()

	return nil
}

func (c *Consumer) Stop() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.running {
		return nil
	}

	c.ticker.Stop()
	close(c.done)
	c.running = false

	return nil
}
