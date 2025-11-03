package try

import (
	"errors"
	"github.com/paveldanilin/go-camel/internal/eip/choice"
	"github.com/paveldanilin/go-camel/internal/eip/setbody"
	"github.com/paveldanilin/go-camel/internal/eip/seterror"
	"github.com/paveldanilin/go-camel/internal/eip/setheader"
	"github.com/paveldanilin/go-camel/internal/expression"
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
	"testing"
)

func TestDoTryProcessor_Error(t *testing.T) {
	var mandatoryParameterMissingErr = errors.New("mandatory parameter missing")

	tryBlock := NewProcessor("", "critical section").
		AddProcessor(setheader.NewProcessor("", "set a", "a", expression.NewConst(1))).
		AddProcessor(setheader.NewProcessor("", "set b", "b", expression.NewConst(1))).
		AddProcessor(setbody.NewProcessor("", "set body", expression.MustSimple("header.a + header.b"))).
		AddProcessor(choice.NewProcessor("", "test body").
			AddWhen(expression.MustSimple("body == 2"), seterror.NewProcessor("", "", mandatoryParameterMissingErr)),
		).
		AddCatch(ErrorContains("mandatory parameter missing"), setheader.NewProcessor("", "", "ERROR", expression.MustSimple("error"))).
		AddFinally(setbody.NewProcessor("", "", expression.NewConst("RESULT")))

	e := exchange.NewExchange(nil)

	tryBlock.Process(e)

	if !errors.Is(e.Message().MustHeader("ERROR").(error), mandatoryParameterMissingErr) {
		t.Errorf("TestDoTryProcessor_Error() = %v; want %v", e.Error(), mandatoryParameterMissingErr)
	}

	expected := "RESULT"
	if e.Message().Body != expected {
		t.Errorf("TestDoTryProcessor_Error() = %v; want %s", e.Message().Body, expected)
	}
}
