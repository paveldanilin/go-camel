package processor

import "github.com/paveldanilin/go-camel/camel"

type Process func(exchange *camel.Exchange)

func (p Process) Process(exchange *camel.Exchange) {
	if !exchange.On("func{}") {
		return
	}

	p(exchange)
}
