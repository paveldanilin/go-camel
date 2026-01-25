package expression

import (
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFunc_Eval(t *testing.T) {
	fnExpr := NewFunc(func(e *exchange.Exchange) (any, error) {
		return e.Message().MustHeader("a").(int) + e.Message().MustHeader("b").(int), nil
	})

	m := exchange.NewExchange(nil)
	m.Message().SetHeader("a", 1)
	m.Message().SetHeader("b", 9)

	result, err := fnExpr.Eval(m)

	assert.NoError(t, err)
	assert.Equal(t, 10, result)
}
