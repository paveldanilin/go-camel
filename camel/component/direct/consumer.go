package direct

import (
	"context"
	"fmt"
	"github.com/paveldanilin/go-camel/camel"
	"sync"
)

type Consumer struct {
	endpoint  *Endpoint
	producers []camel.Producer
	running   bool
	cancel    context.CancelFunc
	mu        sync.Mutex
}

func (c *Consumer) Start() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.running {
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	c.cancel = cancel
	c.running = true

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case message := <-c.endpoint.queue:
				fmt.Printf("%+v\n", len(c.producers))
				for _, producer := range c.producers {
					// TODO: async?
					if err := producer.Process(message); err != nil {
						fmt.Println("Error processing message:", err)
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

	c.cancel()
	c.running = false
	return nil
}
