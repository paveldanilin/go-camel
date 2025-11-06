package direct

import (
	"github.com/paveldanilin/go-camel/pkg/camel/api"
	"github.com/paveldanilin/go-camel/pkg/camel/uri"
	"sync"
)

type Endpoint struct {
	uri      *uri.URI
	mu       sync.RWMutex
	consumer *Consumer
	producer *Producer

	name string
}

func (e *Endpoint) Uri() *uri.URI {
	return e.uri
}

func (e *Endpoint) CreateConsumer(processor api.Processor) (api.Consumer, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.consumer == nil {
		e.consumer = &Consumer{
			endpoint:  e,
			producers: []api.Producer{processor},
		}
	} else {
		e.consumer.producers = append(e.consumer.producers, processor)
	}

	return e.consumer, nil
}

func (e *Endpoint) CreateProducer() (api.Producer, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.producer == nil {
		e.producer = &Producer{endpoint: e}
	}

	return e.producer, nil
}
