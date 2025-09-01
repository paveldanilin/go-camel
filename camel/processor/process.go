package processor

import "github.com/paveldanilin/go-camel/camel"

type Process func(exchange *camel.Exchange)

func (p Process) Process(exchange *camel.Exchange) {
	if err := exchange.CheckCancelOrTimeout(); err != nil {
		exchange.Error = err
		return
	}

	p(exchange)
}
