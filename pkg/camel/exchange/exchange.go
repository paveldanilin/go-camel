package exchange

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type Exchange struct {
	id    string
	start time.Time

	properties Properties
	variables  Variables
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
		properties: Properties{},
		variables:  Variables{},
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

func (e *Exchange) Properties() *Properties {
	return &e.properties
}

func (e *Exchange) Property(name string) (any, bool) {
	return e.properties.Get(name)
}

func (e *Exchange) PropertyOrDefault(name string, def any) any {
	v, exists := e.properties.Get(name)
	if !exists {
		v = def
	}
	return v
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

func (e *Exchange) Variables() *Variables {
	return &e.variables
}

func (e *Exchange) Variable(name string) (any, bool) {
	return e.variables.Get(name)
}

func (e *Exchange) SetVariable(name string, value any) {
	e.variables.Set(name, value)
}

func (e *Exchange) HasVariable(name string) bool {
	return e.variables.Has(name)
}

func (e *Exchange) RemoveVariable(name string) {
	e.variables.Remove(name)
}

func (e *Exchange) Copy() *Exchange {
	if e == nil {
		return nil
	}

	var propsCopy Properties
	if e.properties != nil {
		propsCopy = e.properties.Copy()
	}

	var varsCopy Variables
	if e.variables != nil {
		varsCopy = e.variables.Copy()
	}

	var msgCopy *Message
	if e.message != nil {
		msgCopy = e.message.Copy()
	}

	return &Exchange{
		id:          uuid.NewString(),
		properties:  propsCopy,
		variables:   varsCopy,
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
//	property	- Exchange properties map (k-v)
//	variable	- Exchange variables map (k-v)
func (e *Exchange) AsMap() map[string]any {
	var messageHeaders map[string]any
	if e.message == nil || e.message.headers == nil {
		messageHeaders = map[string]any{}
	} else {
		messageHeaders = e.message.headers
	}

	var properties map[string]any
	if e.properties == nil {
		properties = map[string]any{}
	} else {
		properties = e.properties
	}

	var variables map[string]any
	if e.variables == nil {
		variables = map[string]any{}
	} else {
		variables = e.variables
	}

	return map[string]any{
		"id":         e.Message().Id(),
		"exchangeId": e.Id(),
		"property":   properties,
		"variable":   variables,
		"error":      e.Error(),
		"body":       e.Message().Body(),
		"header":     messageHeaders,
	}
}
