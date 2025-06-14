package main

import (
	"fmt"
	"github.com/paveldanilin/go-camel/camel"
	"github.com/paveldanilin/go-camel/camel/component/direct"
	"github.com/paveldanilin/go-camel/camel/component/timer"
	"github.com/paveldanilin/go-camel/camel/expr"
	"github.com/paveldanilin/go-camel/camel/processor"
	"time"
)

func main() {

	camelContext := camel.NewContext()
	camelContext.RegisterComponent(direct.NewComponent())
	camelContext.RegisterComponent(timer.NewComponent())

	camelContext.RegisterRoute(camel.NewRoute("sum", "direct:sum",
		processor.Pipeline(
			processor.SetPayload(expr.Func(func(message *camel.Message) (any, error) {
				return message.Payload().(int) + 12, nil
			})),
			processor.To("direct:print")),
	))
	camelContext.RegisterRoute(camel.NewRoute("t", "timer:zzz", processor.Pipeline(
		processor.Process(func(message *camel.Message) error {
			message.SetPayload(fmt.Sprintf("COUNT: %d", message.Payload()))
			return nil
		}),
		processor.To("direct:print"))))
	camelContext.RegisterRoute(camel.NewRoute("print", "direct:print", processor.LogMessage("print-1")))
	camelContext.RegisterRoute(camel.NewRoute("print1", "direct:print", processor.LogMessage("print-2")))

	camelContext.Start()

	camelContext.Send("direct:sum", 2, nil)

	time.Sleep(1 * time.Minute)

	camelContext.Stop()
	println("stop")
}
