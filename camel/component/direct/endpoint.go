package direct

import (
	"github.com/paveldanilin/go-camel/camel"
	"sync"
)

type Endpoint struct {
	uri      string
	queue    chan *camel.Message
	mu       sync.RWMutex
	consumer *Consumer
}

func (e *Endpoint) Uri() string {
	return e.uri
}

func (e *Endpoint) Consumer(processor camel.Processor) (camel.Consumer, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.consumer == nil {
		e.consumer = &Consumer{
			endpoint:  e,
			producers: []camel.Producer{processor},
			running:   false,
		}
	} else {
		e.consumer.producers = append(e.consumer.producers, processor)
	}

	return e.consumer, nil
}

func (e *Endpoint) Producer() (camel.Producer, error) {
	return &Producer{endpoint: e}, nil
}

func (e *Endpoint) SendMessage(message *camel.Message) error {

	for _, producer := range e.consumer.producers {
		producer.Process(message)
	}

	return nil
}
