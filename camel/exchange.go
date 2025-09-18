package camel

import (
	"context"
	"github.com/google/uuid"
	"strings"
	"time"
)

type RuntimeProvider interface {
	ExchangeFactory

	Component(componentId string) Component
	Endpoint(uri string) Endpoint
	Route(routeId string) *route
}

type ExchangeFactory interface {
	NewExchange(c context.Context) *Exchange
}

type Exchange struct {
	id      string
	runtime RuntimeProvider
	start   time.Time

	properties Map
	message    *Message
	err        error

	ctx         context.Context
	cancel      context.CancelFunc
	hasDeadline bool
	deadline    time.Time

	// TODO: move to property MessageHistory
	path []string
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
		path:       []string{},
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
	return e.err != nil
}

func (e *Exchange) Error() error {
	return e.err
}

func (e *Exchange) SetError(err error) {
	e.err = err
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

	newPath := make([]string, len(e.path))
	for i, v := range e.path {
		newPath[i] = v
	}

	return &Exchange{
		id:         uuid.NewString(),
		runtime:    e.runtime,
		properties: propsCopy,
		start:      e.start,
		message:    msgCopy,
		err:        e.err,

		ctx:         ctx,
		cancel:      cancel,
		hasDeadline: e.hasDeadline,
		deadline:    e.deadline,

		path: newPath,
	}
}

func (e *Exchange) pushStep(stepName string) {
	e.path = append(e.path, strings.TrimSpace(stepName))
}

// Path returns the stack of steps taken
func (e *Exchange) Path() []string {
	return e.path
}

// On intended to be called at the start of each processor, makes several things:
//  1. pushes step name in stack of steps
//  2. checks context cancellation or timeout
//  3. notifies subscribers
//
// Returns TRUE if processor can proceed with handling this exchange, FALSE - otherwise.
func (e *Exchange) On(stepName string) bool {
	e.pushStep(stepName)

	if err := e.CheckCancelOrTimeout(); err != nil {
		e.SetError(err)
		return false
	}

	return true
}
