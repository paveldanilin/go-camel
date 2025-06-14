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

func (e *Endpoint) Consumer(processor camel.Processor) (camel.Consumer, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.consumer == nil {
		e.consumer = &Consumer{
			processors: []camel.Processor{processor},
			component:  e.component,
		}
	} else {
		e.consumer.processors = append(e.consumer.processors, processor)
	}

	return e.consumer, nil
}

func (e *Endpoint) Producer() (camel.Producer, error) {
	return nil, errors.New("timer: producer not supported")
}
