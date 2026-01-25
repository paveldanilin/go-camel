package processor

import (
	"context"
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
	"github.com/stretchr/testify/assert"
	"testing"
)

func mulBody(n int) ExpressionFunc {
	return func(e *exchange.Exchange) (any, error) {
		v := e.MessageBody().(int)
		return v * n, nil
	}
}

func TestMulticast_Sync(t *testing.T) {
	multicast := NewMulticast(false, true, nil)

	multicast.AddProcessor(NewSetBody(mulBody(2)))
	multicast.AddProcessor(NewSetBody(mulBody(3)))
	multicast.AddProcessor(NewSetBody(mulBody(4)))
	multicast.AddProcessor(NewSetBody(mulBody(5)))

	ex := exchange.NewExchange(context.TODO()).SetMessageBody(2)
	originalExchangeId := ex.Id()

	multicast.Process(ex)

	assert.NoError(t, ex.Error(), "Exchange contains error")
	assert.Equal(t, originalExchangeId, ex.Id(),
		"Sync processing without aggregator must keep original Exchange and discard all results from processors")
	assert.Equal(t, 2, ex.MessageBody(),
		"Sync processing without aggregator must keep original Exchange and discard all results from processors")
}

func TestMulticast_SyncUseLatestAgg(t *testing.T) {
	multicast := NewMulticast(false, true, exchange.UseLatestAgg())

	multicast.AddProcessor(NewSetBody(mulBody(2)))
	multicast.AddProcessor(NewSetBody(mulBody(3)))
	multicast.AddProcessor(NewSetBody(mulBody(4)))
	multicast.AddProcessor(NewSetBody(mulBody(5)))

	ex := exchange.NewExchange(context.TODO()).SetMessageBody(2)
	originalExchangeId := ex.Id()
	originalBody := 2

	multicast.Process(ex)

	assert.NoError(t, ex.Error(), "Exchange contains error")
	assert.NotEqualf(t, originalExchangeId, ex.Id(),
		"Sync processing with UseLatestAgg must returns the last Exchange")
	assert.NotEqual(t, originalBody, ex.MessageBody())
	assert.Equal(t, 10, ex.MessageBody(),
		"Sync processing with UseLatestAgg must returns the last Exchange")
}

func TestMulticast_SyncGroupedBodyAgg(t *testing.T) {
	multicast := NewMulticast(false, true, exchange.GroupedBodyAgg())

	multicast.AddProcessor(NewSetBody(constant("service1: result")))
	multicast.AddProcessor(NewSetBody(constant("service2: result")))
	multicast.AddProcessor(NewSetBody(constant("service3: result")))
	multicast.AddProcessor(NewSetBody(constant("service4: result")))

	ex := exchange.NewExchange(context.TODO()).SetMessageBody("call services: 1,2,3,4")

	multicast.Process(ex)

	assert.NoError(t, ex.Error(), "Exchange contains error")
	assert.Equal(t,
		[]any{"service1: result", "service2: result", "service3: result", "service4: result"},
		ex.MessageBody(),
		"Sync processing with GroupedBodyAgg must aggregate all output results Exchange's body")
}

func TestMulticast_Parallel(t *testing.T) {
	multicast := NewMulticast(true, true, exchange.GroupedBodyAgg())

	multicast.AddProcessor(
		NewPipeline(true).
			AddProcessor(NewDelay(500)).
			AddProcessor(NewSetBody(constant("service1: result"))),
	)
	multicast.AddProcessor(
		NewPipeline(true).
			AddProcessor(NewDelay(2000)).
			AddProcessor(NewSetBody(constant("service2: slow result"))),
	)
	multicast.AddProcessor(
		NewPipeline(true).
			AddProcessor(NewDelay(100)).
			AddProcessor(NewSetBody(constant("service3: result"))),
	)
	multicast.AddProcessor(
		NewPipeline(true).
			AddProcessor(NewDelay(100)).
			AddProcessor(NewSetBody(constant("service4: result"))),
	)

	ex := exchange.NewExchange(context.TODO()).SetMessageBody("call services: 1,2,3,4")

	multicast.Process(ex)

	assert.NoError(t, ex.Error(), "Exchange contains error")
	// We get the same order of results due to parallel processing approach
	assert.Equal(t,
		[]any{"service1: result", "service2: slow result", "service3: result", "service4: result"},
		ex.MessageBody(),
		"Sync processing with GroupedBodyAgg must aggregate all output results Exchange's body")
}
