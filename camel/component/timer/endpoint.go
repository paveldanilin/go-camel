package timer

import (
	"errors"
	"fmt"
	"github.com/paveldanilin/go-camel/camel"
	"github.com/paveldanilin/go-camel/camel/uri"
	"sync"
	"time"
)

const EndpointParamInterval = "interval"

type Endpoint struct {
	uri       *uri.URI
	mu        sync.RWMutex
	component *Component
	consumer  *Consumer

	name     string
	interval time.Duration
}

func NewEndpoint(uri *uri.URI, c *Component) (*Endpoint, error) {
	interval, err := resolveInterval(uri)
	if err != nil {
		return nil, err
	}

	timerEndpoint := &Endpoint{
		component: c,
		uri:       uri,
		name:      uri.Path(),
		interval:  interval,
	}

	return timerEndpoint, nil
}

func resolveInterval(uri *uri.URI) (time.Duration, error) {
	if !uri.HasParam(EndpointParamInterval) {
		return 0, fmt.Errorf("timer: mandatory parameter not found '%s'", EndpointParamInterval)
	}

	return time.ParseDuration(uri.MustParam(EndpointParamInterval))
}

func (e *Endpoint) Uri() *uri.URI {
	return e.uri
}

func (e *Endpoint) CreateConsumer(processor camel.Processor) (camel.Consumer, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.consumer == nil {
		consumer, err := NewConsumer(e)
		if err != nil {
			return nil, err
		}

		consumer.processors = append(consumer.processors, processor)
		e.consumer = consumer
	} else {
		e.consumer.processors = append(e.consumer.processors, processor)
	}

	return e.consumer, nil
}

func (e *Endpoint) CreateProducer() (camel.Producer, error) {
	return nil, errors.New("timer: producer not supported")
}
