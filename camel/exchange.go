package camel

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type RuntimeProvider interface {
	ExchangeFactory

	Component(componentId string) Component
	Endpoint(uri string) Endpoint
	Route(routeId string) *Route
}

type ExchangeFactory interface {
	NewExchange(c context.Context) *Exchange
}

type Exchange struct {
	id         string
	runtime    RuntimeProvider
	properties Map
	start      time.Time
	message    *Message
	Error      error

	ctx         context.Context
	cancel      context.CancelFunc
	hasDeadline bool
	deadline    time.Time
}

func NewExchange(c context.Context, r RuntimeProvider) *Exchange {
	if c == nil {
		c = context.Background()
	}
	ctx, cancel := context.WithCancel(c)

	e := &Exchange{
		id:         uuid.NewString(),
		runtime:    r,
		start:      time.Now(),
		properties: newMap(),
		message:    NewMessage(),
		ctx:        ctx,
		cancel:     cancel,
	}

	if dl, ok := c.Deadline(); ok {
		e.hasDeadline = true
		e.deadline = dl
	}
	return e
}

func (e *Exchange) Id() string {
	return e.id
}

func (e *Exchange) Context() context.Context {
	return e.ctx
}

func (e *Exchange) Cancel() bool {
	if e.cancel != nil {
		e.cancel()
		return true
	}
	return false
}

func (e *Exchange) Deadline() (time.Time, bool) {
	return e.deadline, e.hasDeadline
}

func (e *Exchange) HasDeadline() bool {
	return e.hasDeadline
}

func (e *Exchange) DeadlineExceeded() bool {
	return e.hasDeadline && time.Now().After(e.deadline)
}

func (e *Exchange) WaitOrErr() error {
	<-e.Context().Done()
	return e.Context().Err()
}

func (e *Exchange) CheckCancelOrTimeout() error {
	if e.DeadlineExceeded() {
		// Cheap test
		return context.DeadlineExceeded
	}

	select {
	case <-e.Context().Done():
		return e.Context().Err()
	default:
		return nil
	}
}

func (e *Exchange) Runtime() RuntimeProvider {
	return e.runtime
}

func (e *Exchange) Properties() *Map {
	return &e.properties
}

func (e *Exchange) Property(name string) (any, bool) {
	return e.properties.Get(name)
}

func (e *Exchange) SetProperty(name string, value any) {
	e.properties.Set(name, value)
}

func (e *Exchange) Message() *Message {
	return e.message
}

func (e *Exchange) StartedAt() time.Time {
	return e.start
}

func (e *Exchange) IsError() bool {
	return e.Error != nil
}

func (e *Exchange) Copy() *Exchange {
	if e == nil {
		return nil
	}

	var propsCopy Map
	if e.properties != nil {
		propsCopy = e.properties.Copy()
	}
	var msgCopy *Message
	if e.message != nil {
		msgCopy = e.message.Copy()
	}

	// Inherit parent context
	var (
		ctx    context.Context
		cancel context.CancelFunc
	)
	if e.hasDeadline && !e.deadline.IsZero() {
		// Inherit from parent
		ctx, cancel = context.WithDeadline(e.ctx, e.deadline)
	} else {
		ctx, cancel = context.WithCancel(e.ctx)
	}

	return &Exchange{
		id:         uuid.NewString(),
		runtime:    e.runtime,
		properties: propsCopy,
		start:      e.start,
		message:    msgCopy,
		Error:      e.Error,

		ctx:         ctx,
		cancel:      cancel,
		hasDeadline: e.hasDeadline,
		deadline:    e.deadline,
	}
}