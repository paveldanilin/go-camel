package direct

import (
	"github.com/paveldanilin/go-camel/camel"
	"sync"
)

type Endpoint struct {
	uri      string
	mu       sync.RWMutex
	consumer *Consumer
	producer *Producer
}

func (e *Endpoint) Uri() string {

	return e.uri
}

func (e *Endpoint) CreateConsumer(processor camel.Processor) (camel.Consumer, error) {

	e.mu.Lock()
	defer e.mu.Unlock()

	if e.consumer == nil {
		e.consumer = &Consumer{
			endpoint:  e,
			producers: []camel.Producer{processor},
		}
	} else {
		e.consumer.producers = append(e.consumer.producers, processor)
	}

	return e.consumer, nil
}

func (e *Endpoint) CreateProducer() (camel.Producer, error) {

	e.mu.Lock()
	defer e.mu.Unlock()

	if e.producer == nil {
		e.producer = &Producer{endpoint: e}
	}

	return e.producer, nil
}
