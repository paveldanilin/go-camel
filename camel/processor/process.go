package processor

import "github.com/paveldanilin/go-camel/camel"

type Process func(message *camel.Message)

func (p Process) Process(message *camel.Message) {

	p(message)
}
