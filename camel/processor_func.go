package camel

type funcProcessor func(exchange *Exchange)

func (p funcProcessor) Process(exchange *Exchange) {
	if !exchange.On("func{}") {
		return
	}

	p(exchange)
}
