package direct

import (
	"github.com/paveldanilin/go-camel/pkg/camel"
)

type Consumer struct {
	endpoint  *Endpoint
	producers []camel.Producer
}

func (c *Consumer) Start() error {
	return nil
}

func (c *Consumer) Stop() error {
	return nil
}
