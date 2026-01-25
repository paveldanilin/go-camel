package processor

import (
	"context"
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoopCountProcessor(t *testing.T) {
	loop := NewRepeatLoop(5).
		AddProcessor(NewSetBody(ExpressionFunc(func(e *exchange.Exchange) (any, error) {
			it, _ := e.Property(exchange.PropertyLoopIndex)
			return it, nil
		})))

	e := exchange.NewExchange(context.TODO())

	loop.Process(e)

	assert.Equal(t, 4, e.MessageBody())
}
