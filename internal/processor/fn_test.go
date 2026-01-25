package processor

import (
	"context"
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFnProcessor(t *testing.T) {
	p := NewFn(func(e *exchange.Exchange) {
		e.Message().SetBody("Hello, World!")
	})

	e := exchange.NewExchange(context.TODO())
	p.Process(e)

	assert.Equal(t, "Hello, World!", e.MessageBody())
}
