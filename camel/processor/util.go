package processor

import (
	"fmt"
	"github.com/paveldanilin/go-camel/camel"
)

// Invoke invokes processor with recovery
func Invoke(p camel.Processor, message *camel.Message) (panicked bool) {

	defer func() {
		if r := recover(); r != nil {
			message.Error = fmt.Errorf("panic recovered: %v", r)
			panicked = true
		}
	}()

	p.Process(message)
	return false
}
