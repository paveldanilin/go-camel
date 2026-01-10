package processor

import (
	"github.com/paveldanilin/go-camel/internal/expression"
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
	"testing"
)

func TestChoiceProcessor(t *testing.T) {
	c := NewChoice().
		When(expression.MustSimple("header.val > 5"), NewSetBody(expression.NewConstant(555))).
		When(expression.MustSimple("header.val < 5"), NewSetBody(expression.NewConstant(777)))

	e := exchange.NewExchange(nil)
	e.Message().SetHeader("val", 2)

	c.Process(e)

	expected := 777
	if e.Message().Body() != expected {
		t.Errorf("TestChoiceProcessor() = %v; want %v", e.Message().Body(), expected)
	}
}
