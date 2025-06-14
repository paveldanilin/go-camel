package processor

import "github.com/paveldanilin/go-camel/camel"

type Process func(message *camel.Message) error

func (p Process) Process(message *camel.Message) error {

	return p(message)
}
