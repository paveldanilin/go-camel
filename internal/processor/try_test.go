package processor

import (
	"errors"
	"github.com/paveldanilin/go-camel/internal/expression"
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
	"testing"
)

func TestDoTryProcessor_Error(t *testing.T) {
	var mandatoryParameterMissingErr = errors.New("mandatory parameter missing")

	tryBlock := NewTry().
		AddProcessor(NewSetHeader("a", expression.NewConstant(1))).
		AddProcessor(NewSetHeader("b", expression.NewConstant(1))).
		AddProcessor(NewSetBody(expression.MustSimple("header.a + header.b"))).
		AddProcessor(NewChoice().
			When(expression.MustSimple("body == 2"), NewSetError(mandatoryParameterMissingErr)),
		).
		AddCatch(ErrorContains("mandatory parameter missing"), NewSetHeader("ERROR", expression.MustSimple("error"))).
		AddFinally(NewSetBody(expression.NewConstant("RESULT")))

	e := exchange.NewExchange(nil)

	tryBlock.Process(e)

	if !errors.Is(e.Message().MustHeader("ERROR").(error), mandatoryParameterMissingErr) {
		t.Errorf("TestDoTryProcessor_Error() = %v; want %v", e.Error(), mandatoryParameterMissingErr)
	}

	expected := "RESULT"
	if e.Message().Body() != expected {
		t.Errorf("TestDoTryProcessor_Error() = %v; want %s", e.Message().Body(), expected)
	}
}
