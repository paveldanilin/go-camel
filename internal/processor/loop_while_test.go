package processor

import (
	"context"
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWhileLoopProcessor(t *testing.T) {
	loop, err := NewWhileLoop(PredicateFunc(func(e *exchange.Exchange) (bool, error) {
		it, _ := e.Property(exchange.PropertyLoopIndex)
		return it.(int) < 10, nil
	}), true)
	assert.NoError(t, err)

	loop.AddProcessor(P(func(e *exchange.Exchange) {
		e.SetMessageBody(e.PropertyOrDefault(exchange.PropertyLoopIndex, 0))
	}))

	e := exchange.NewExchange(context.TODO())

	loop.Process(e)

	assert.Equal(t, 9, e.MessageBody())
}

func TestWhileLoopProcessor_Break(t *testing.T) {
	loop, err := NewWhileLoop(PredicateFunc(func(e *exchange.Exchange) (bool, error) {
		return true, nil
	}), true)
	assert.NoError(t, err)

	loop.AddProcessor(P(func(e *exchange.Exchange) {
		e.SetMessageBody(e.PropertyOrDefault(exchange.PropertyLoopIndex, 0))
	}))

	breakCond := NewChoice()
	// break the loop after 5 iterations
	err = breakCond.AddWhen(
		PredicateFunc(func(e *exchange.Exchange) (bool, error) {
			return e.MessageBody().(int) >= 5, nil
		}),
		NewSetProperty(exchange.PropertyLoopBreak, constant(true)))
	assert.NoError(t, err)

	loop.AddProcessor(breakCond)

	e := exchange.NewExchange(nil)

	loop.Process(e)

	assert.Equal(t, 5, e.MessageBody())
}
