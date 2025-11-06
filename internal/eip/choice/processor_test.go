package choice

import (
	"github.com/paveldanilin/go-camel/internal/eip/setbody"
	"github.com/paveldanilin/go-camel/internal/expression"
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
	"testing"
)

func TestChoiceProcessor(t *testing.T) {
	c := NewProcessor("", "test val").
		AddWhen(expression.MustSimple("header.val > 5"), setbody.NewProcessor("", "", expression.NewConst(555))).
		AddWhen(expression.MustSimple("header.val < 5"), setbody.NewProcessor("", "", expression.NewConst(777)))

	e := exchange.NewExchange(nil)
	e.Message().SetHeader("val", 2)

	c.Process(e)

	expected := 777
	if e.Message().Body != expected {
		t.Errorf("TestChoiceProcessor() = %v; want %v", e.Message().Body, expected)
	}
}
