package exchange

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type Exchange struct {
	id    string
	start time.Time

	properties Map
	message    *Message
	err        error

	ctx         context.Context
	cancel      context.CancelFunc
	hasDeadline bool
	deadline    time.Time
}

func NewExchange(c context.Context) *Exchange {
	if c == nil {
		c = context.Background()
	}
	ctx, cancel := context.WithCancel(c)

	e := &Exchange{
		id:         uuid.NewString(),
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

func (e *Exchange) Properties() *Map {
	return &e.properties
}

func (e *Exchange) Property(name string) (any, bool) {
	return e.properties.Get(name)
}

func (e *Exchange) SetProperty(name string, value any) {
	e.properties.Set(name, value)
}

func (e *Exchange) HasProperty(name string) bool {
	return e.properties.Has(name)
}

func (e *Exchange) RemoveProperty(name string) {
	e.properties.Remove(name)
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

	return &Exchange{
		id:          uuid.NewString(),
		properties:  propsCopy,
		start:       e.start,
		message:     msgCopy,
		err:         e.err,
		ctx:         e.ctx,
		cancel:      e.cancel,
		hasDeadline: e.hasDeadline,
		deadline:    e.deadline,
	}
}

// AsMap returns Exchange's data as map.
//
// Keys:
//
//	id 			- Message is
//	exchangeId 	- Exchange id
//	body 		- Message body
//	header 		- Message headers map (k-v)
//	error		- Exchange error
//	property	- Exchange properties map (kv-)
func (e *Exchange) AsMap() map[string]any {
	return map[string]any{
		"id":         e.Message().Id(),
		"exchangeId": e.Id(),
		"body":       e.Message().Body,
		"header":     e.Message().Headers().All(),
		"error":      e.Error(),
		"property":   e.Properties().All(),
	}
}
