package exchange

import "fmt"

type AggregatorFunc func(oldExchange, newExchange *Exchange) *Exchange

func (f AggregatorFunc) AggregateExchange(oldExchange, newExchange *Exchange) *Exchange {
	return f(oldExchange, newExchange)
}

// GroupedBodyAgg implements 'Collect result body' aggregation strategy.
// Wraps Exchange bodies as slice.
// Keeps first Exchange data (properties, variables, message headers).
func GroupedBodyAgg() AggregatorFunc {
	return func(oldExchange, newExchange *Exchange) *Exchange {
		if oldExchange == nil {
			newExchange.SetMessageBody([]any{newExchange.Message().Body()})
			return newExchange
		}

		if list, isList := oldExchange.Message().Body().([]any); isList {
			return oldExchange.SetMessageBody(append(list, newExchange.Message().Body()))
		}

		return oldExchange
	}
}

// UseLatestAgg implements 'Last result' aggregation strategy, returns the last Exchange.
func UseLatestAgg() AggregatorFunc {
	return func(oldExchange, newExchange *Exchange) *Exchange {
		return newExchange
	}
}

// StringConcatAgg implements 'Concat result body as string' aggregation strategy.
// Keeps first Exchange data (properties, variables, message headers).
func StringConcatAgg(delimiter string) AggregatorFunc {
	return func(oldExchange, newExchange *Exchange) *Exchange {
		if oldExchange == nil {
			//
			return newExchange
		}
		return oldExchange.SetMessageBody(
			fmt.Sprintf("%v%s%v", oldExchange.MessageBody(), delimiter, newExchange.MessageBody()),
		)
	}
}
