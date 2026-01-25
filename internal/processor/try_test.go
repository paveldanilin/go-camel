package processor

import (
	"context"
	"errors"
	"github.com/paveldanilin/go-camel/internal/expression"
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDoTryProcessor_Error(t *testing.T) {
	var mandatoryParameterMissingErr = errors.New("mandatory parameter missing")

	setErrCond := NewChoice()
	err := setErrCond.AddWhen(PredicateFunc(func(e *exchange.Exchange) (bool, error) {
		return e.MessageBody().(int) == 2, nil
	}), NewSetError(mandatoryParameterMissingErr))
	assert.NoError(t, err)

	tryBlock := NewTry().
		AddProcessor(NewSetHeader("a", expression.NewConstant(1))).
		AddProcessor(NewSetHeader("b", expression.NewConstant(1))).
		AddProcessor(NewSetBody(expression.MustSimple("header.a + header.b"))).
		AddProcessor(setErrCond).
		AddCatch(ErrorContains("mandatory parameter missing"), NewSetHeader("ERROR", expression.MustSimple("error"))).
		AddFinally(NewSetBody(expression.NewConstant("RESULT")))

	e := exchange.NewExchange(context.TODO())

	tryBlock.Process(e)

	//if !errors.Is(e.Message().MustHeader("ERROR").(error), mandatoryParameterMissingErr) {
	//	t.Errorf("TestDoTryProcessor_Error() = %v; want %v", e.Error(), mandatoryParameterMissingErr)
	//}
	assert.EqualError(t, mandatoryParameterMissingErr, e.Message().MustHeader("ERROR").(error).Error())

	assert.Equal(t, "RESULT", e.MessageBody())
}
