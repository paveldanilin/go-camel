package timer

import (
	"fmt"
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
	component  *Component
	processors []camel.Processor

	ticker   *time.Ticker
	interval time.Duration
}

func NewConsumer(params map[string]any) (*Consumer, error) {

	dur, err := time.ParseDuration(params["interval"].(string))
	if err != nil {
		return nil, err
	}

	return &Consumer{
		processors: []camel.Processor{},
		interval:   dur,
	}, nil
}

func (c *Consumer) Start() error {

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.running {
		return nil
	}

	c.ticker = time.NewTicker(c.interval)
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

				message := camel.NewMessageWithContext(c.component.runtime)

				message.SetHeader(HeaderTimeFiredTime, t)
				message.SetHeader(HeaderTimerName, "")
				message.SetHeader(HeaderTimerCounter, count)

				for _, processor := range c.processors {
					if err := processor.Process(message); err != nil {
						fmt.Println("Error processing timer exchange:", err)
					}
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
