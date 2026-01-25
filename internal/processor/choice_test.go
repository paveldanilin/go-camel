package processor

import (
	"context"
	"errors"
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestChoiceProcessor(t *testing.T) {

	prdHeaderGt := func(headerName string, limit int) Predicate {
		return PredicateFunc(func(e *exchange.Exchange) (bool, error) {
			if v, exists := e.MessageHeader(headerName); exists {
				return v.(int) > limit, nil
			}
			return false, errors.New("header not found: " + headerName)
		})
	}

	prdHeaderLt := func(headerName string, limit int) Predicate {
		return PredicateFunc(func(e *exchange.Exchange) (bool, error) {
			if v, exists := e.MessageHeader(headerName); exists {
				return v.(int) < limit, nil
			}
			return false, errors.New("header not found: " + headerName)
		})
	}

	choice := NewChoice()

	err := choice.AddWhen(prdHeaderGt("val", 5), P(func(e *exchange.Exchange) {
		e.SetMessageBody(555)
	}))
	assert.NoError(t, err)

	err = choice.AddWhen(prdHeaderLt("val", 5), P(func(e *exchange.Exchange) {
		e.SetMessageBody(777)
	}))
	assert.NoError(t, err)

	// Test 1: 777
	e := exchange.NewExchange(context.TODO()).SetMessageHeader("val", 2)
	choice.Process(e)
	assert.Equal(t, 777, e.MessageBody())

	// Test 2: 555
	e.SetMessageHeader("val", 10)
	choice.Process(e)
	assert.Equal(t, 555, e.MessageBody())
}
