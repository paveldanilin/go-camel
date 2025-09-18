package camel

import "testing"

func TestChoiceProcessor(t *testing.T) {
	choice := newChoiceProcessor().
		When(mustSimpleExpr("header.val > 5"), newSetBodyProcessor(newConstExpr(555))).
		When(mustSimpleExpr("header.val < 5"), newSetBodyProcessor(newConstExpr(777)))

	exchange := NewExchange(nil, nil)
	exchange.Message().SetHeader("val", 2)

	choice.Process(exchange)

	expected := 777
	if exchange.Message().Body != expected {
		t.Errorf("TestChoiceProcessor() = %v; want %v", exchange.Message().Body, expected)
	}
}
