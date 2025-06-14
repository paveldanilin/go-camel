package timer

import (
	"fmt"
	"github.com/paveldanilin/go-camel/camel"
	"sync"
	"time"
)

type Consumer struct {
	component  *Component
	processors []camel.Processor
	ticker     *time.Ticker
	done       chan struct{}
	running    bool
	mu         sync.Mutex
}

func (c *Consumer) Start() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.running {
		return nil
	}

	c.ticker = time.NewTicker(2 * time.Second)
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
				/*
					message := &camel.Message{
						Body: count,
						Headers: map[string]interface{}{
							"firedTime": t,
						},
						Properties: map[string]interface{}{
							"CamelTimerName":    c.endpoint.name,
							"CamelTimerCounter": count,
						},
					}*/
				message := camel.NewMessageWithContext(c.component.context)

				message.SetPayload(count)
				message.SetHeader("firedTime", t)
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
