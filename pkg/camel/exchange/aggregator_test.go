package exchange

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAggregatorFunc_GroupedBodyAgg(t *testing.T) {
	agg := GroupedBodyAgg()
	exs := []*Exchange{
		NewExchange(context.TODO()).
			SetMessageBody(1).
			SetProperty(PropertyStepName, "Step1"), // GroupedBodyAgg keeps first exchange data
		NewExchange(context.TODO()).
			SetMessageBody(2).
			SetProperty(PropertyStepName, "Step2"),
		NewExchange(context.TODO()).
			SetMessageBody(3).
			SetProperty(PropertyStepName, "Step3"),
		NewExchange(context.TODO()).
			SetMessageBody(4).
			SetProperty(PropertyStepName, "Step4"),
	}
	var result *Exchange = nil

	for _, ex := range exs {
		result = agg.AggregateExchange(result, ex)
	}

	// GroupedBodyAgg keeps first exchange data (properties, variables, message headers)
	wantStepName := "Step1"
	assert.Equalf(t, wantStepName, StepName(result), "Expected step name is '%s'", wantStepName)

	wantBody := []any{1, 2, 3, 4}
	assert.Equalf(t, wantBody, result.MessageBody(), "Expected message body is [1,2,3,4]")
}

func TestAggregatorFunc_UseLatestAgg(t *testing.T) {
	agg := UseLatestAgg()
	exs := []*Exchange{
		NewExchange(context.TODO()).
			SetMessageBody(1).
			SetProperty(PropertyStepName, "Step1"),
		NewExchange(context.TODO()).
			SetMessageBody(2).
			SetProperty(PropertyStepName, "Step2"),
		NewExchange(context.TODO()).
			SetMessageBody(3).
			SetProperty(PropertyStepName, "Step3"),
		NewExchange(context.TODO()).
			SetMessageBody(4).
			SetProperty(PropertyStepName, "Step4"),
	}
	var result *Exchange = nil

	for _, ex := range exs {
		result = agg.AggregateExchange(result, ex)
	}

	wantStepName := "Step4"
	assert.Equalf(t, wantStepName, StepName(result), "Expected step name is '%s'", wantStepName)

	wantBody := 4
	assert.Equalf(t, wantBody, result.MessageBody(), "Expected message body is 4")
}

func TestAggregatorFunc_StringConcatAgg(t *testing.T) {
	agg := StringConcatAgg(";")
	exs := []*Exchange{
		NewExchange(context.TODO()).
			SetMessageBody(1).
			SetProperty(PropertyStepName, "Step1"),
		NewExchange(context.TODO()).
			SetMessageBody("2").
			SetProperty(PropertyStepName, "Step2"),
		NewExchange(context.TODO()).
			SetMessageBody(true).
			SetProperty(PropertyStepName, "Step3"),
		NewExchange(context.TODO()).
			SetMessageBody([]int{5, 5, 5}).
			SetProperty(PropertyStepName, "Step4"),
	}
	var result *Exchange = nil

	for _, ex := range exs {
		result = agg.AggregateExchange(result, ex)
	}

	wantStepName := "Step1"
	assert.Equalf(t, wantStepName, StepName(result), "Expected step name is '%s'", wantStepName)

	wantBody := "1;2;true;[5 5 5]"
	assert.Equalf(t, wantBody, result.MessageBody(), "Expected message body is 4")
}
