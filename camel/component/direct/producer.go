package direct

import (
	"fmt"
	"github.com/paveldanilin/go-camel/camel"
)

type Producer struct {
	endpoint *Endpoint
}

func (p *Producer) Process(message *camel.Message) error {
	println(">>>>>>>")
	select {
	case p.endpoint.queue <- message:
		return nil
	default:
		return fmt.Errorf("queue is full for endpoint %s", p.endpoint.uri)
	}
}
