package camel

import (
	"errors"
	"testing"
)

func TestDoTryProcessor_Error(t *testing.T) {
	var mandatoryParameterMissingErr = errors.New("mandatory parameter missing")

	tryBlock := newTryProcessor("critical section").
		addProcessor(newSetHeaderProcessor("set a", "a", newConstExpression(1))).
		addProcessor(newSetHeaderProcessor("set b", "b", newConstExpression(1))).
		addProcessor(newSetBodyProcessor("set body", mustSimpleExpression("headers.a + headers.b"))).
		addProcessor(newChoiceProcessor("test body").
			addWhen(mustSimpleExpression("body == 2"), newSetErrorProcessor("", mandatoryParameterMissingErr)),
		).
		addCatch(errorContains("mandatory parameter missing"), newSetHeaderProcessor("", "ERROR", mustSimpleExpression("error"))).
		addFinally(newSetBodyProcessor("", newConstExpression("RESULT")))

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
