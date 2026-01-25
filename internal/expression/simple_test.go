package expression

import (
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSimple_Eval(t *testing.T) {
	e, err := NewSimple("header.a == 1")
	assert.NoError(t, err)

	m := exchange.NewExchange(nil)
	m.Message().SetHeader("a", 1)

	ret, err := e.Eval(m)

	assert.NoError(t, err)
	assert.True(t, ret.(bool))
}
