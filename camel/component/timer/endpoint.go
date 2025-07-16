package timer

import (
	"errors"
	"github.com/paveldanilin/go-camel/camel"
	"sync"
)

type Endpoint struct {
	uri       string
	mu        sync.RWMutex
	component *Component
	consumer  *Consumer
}

func (e *Endpoint) Uri() string {

	return e.uri
}

func (e *Endpoint) CreateConsumer(processor camel.Processor) (camel.Consumer, error) {

	e.mu.Lock()
	defer e.mu.Unlock()

	if e.consumer == nil {
		consumer, err := NewConsumer(map[string]any{
			"interval": "5s",
		})
		if err != nil {
			return nil, err
		}

		consumer.processors = append(consumer.processors, processor)
		consumer.component = e.component

		e.consumer = consumer
	} else {
		e.consumer.processors = append(e.consumer.processors, processor)
	}

	return e.consumer, nil
}

func (e *Endpoint) CreateProducer() (camel.Producer, error) {

	return nil, errors.New("timer: producer not supported")
}
