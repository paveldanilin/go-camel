package direct

import (
	"github.com/paveldanilin/go-camel/pkg/camel/api"
)

type Consumer struct {
	endpoint  *Endpoint
	producers []api.Producer
}

func (c *Consumer) Start() error {
	return nil
}

func (c *Consumer) Stop() error {
	return nil
}
