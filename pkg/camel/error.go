package camel

import (
	"errors"
	"fmt"
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
	"time"
)

type Error struct {
	panicValue any
	stackTrace string
	timestamp  time.Time
	exchange   *exchange.Exchange
}

func NewPanicError(ex *exchange.Exchange, panicValue any, stack string) *Error {
	return &Error{
		panicValue: panicValue,
		stackTrace: stack,
		timestamp:  time.Now(),
		exchange:   ex,
	}
}

func (e *Error) PanicValue() any {
	return e.panicValue
}

func (e *Error) IsPanic() bool {
	return e.panicValue != nil
}

func (e *Error) RouteName() string {
	if e.exchange == nil {
		return ""
	}
	return e.exchange.PropertyOrDefault(exchange.PropertyRouteName, "").(string)
}

func (e *Error) RouteStep() string {
	if e.exchange == nil {
		return ""
	}
	return e.exchange.PropertyOrDefault(exchange.PropertyStepName, "").(string)
}

func (e *Error) Exchange() *exchange.Exchange {
	return e.exchange
}

func (e *Error) StackTrace() string {
	return e.stackTrace
}

func (e *Error) Error() string {
	if e.exchange == nil {
		return fmt.Sprintf("step [%s] in route '%s' (exchange %s)", e.RouteStep(), e.RouteName(), "")
	}
	if e.IsPanic() {
		return fmt.Sprintf("step [%s] in route '%s' (exchange %s): %v", e.RouteStep(), e.RouteName(), e.exchange.Id(), e.panicValue)
	}
	return fmt.Sprintf("step [%s] in route '%s' (exchange %s): %v", e.RouteStep(), e.RouteName(), e.exchange.Id(), e.exchange.Error())
}

func (e *Error) Unwrap() error {
	if e.exchange == nil {
		return e.panicValueAsError()
	}
	return e.exchange.Error()
}

func (e *Error) panicValueAsError() error {
	if e.panicValue == nil {
		return nil
	}
	switch t := e.panicValue.(type) {
	case error:
		return t
	case string:
		return errors.New(t)
	default:
		return fmt.Errorf("panic: %v", t)
	}
}
