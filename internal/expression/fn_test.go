package expression

import (
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
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
	if err != nil {
		t.Fatalf("TestFunc_Eval(): %s", err)
	}

	wantResult := 10
	if result != wantResult {
		t.Fatalf("TestFunc_Eval(): expected result %v, but got %v", wantResult, result)
	}
}
