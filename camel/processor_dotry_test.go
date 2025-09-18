package camel

import (
	"errors"
	"testing"
)

func TestDoTryProcessor_Error(t *testing.T) {
	var mandatoryParameterMissingErr = errors.New("mandatory parameter missing")

	tryBlock := newDoTryProcessor(
		newSetHeaderProcessor("a", newConstExpr(1)),
		newSetHeaderProcessor("b", newConstExpr(1)),
		newSetBodyProcessor(mustSimpleExpr("header.a + header.b")),
		newChoiceProcessor().When(mustSimpleExpr("body == 2"), newSetErrorProcessor(mandatoryParameterMissingErr)),
	).
		Catch(ErrorContains("mandatory parameter missing"),
			newSetHeaderProcessor("ERROR", mustSimpleExpr("exchange.error"))).
		Finally(newSetBodyProcessor(newConstExpr("RESULT")))

	exchange := NewExchange(nil, nil)

	tryBlock.Process(exchange)

	if !errors.Is(exchange.Message().MustHeader("ERROR").(error), mandatoryParameterMissingErr) {
		t.Errorf("TestDoTryProcessor_Error() = %v; want %v", exchange.Error(), mandatoryParameterMissingErr)
	}

	expected := "RESULT"
	if exchange.Message().Body != expected {
		t.Errorf("TestDoTryProcessor_Error() = %v; want %s", exchange.Message().Body, expected)
	}
}
