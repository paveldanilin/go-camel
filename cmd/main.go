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

	camelRuntime := camel.NewRuntime()
	camelRuntime.RegisterComponent(direct.NewComponent())
	camelRuntime.RegisterComponent(timer.NewComponent())

	camelRuntime.RegisterRoute(camel.NewRoute("sum", "direct:sum",
		processor.Pipeline(
			processor.SetPayload(expr.Func(func(message *camel.Message) (any, error) {
				return message.MustHeader("a").(int) + message.MustHeader("b").(int), nil
			})),
			processor.To("direct:print"),
			processor.Choice().
				When(expr.MustSimple("header.a == 1"), processor.Process(func(message *camel.Message) error {
					message.SetPayload(444)
					return nil
				})).
				Otherwise(processor.To("direct:x")),
		),
	))

	camelRuntime.RegisterRoute(camel.NewRoute("t", "timer:zzz", processor.Pipeline(
		processor.Process(func(message *camel.Message) error {
			message.SetPayload(fmt.Sprintf("COUNT: %d", message.Payload()))
			return nil
		}),
		processor.To("direct:print"))))

	camelRuntime.RegisterRoute(camel.NewRoute("print", "direct:print", processor.LogMessage("print-1")))

	camelRuntime.RegisterRoute(camel.NewRoute("print1", "direct:print", processor.LogMessage("print-2")))

	camelRuntime.RegisterRoute(camel.NewRoute("x", "direct:x", processor.LogMessage("xxx")))

	camelRuntime.Start()

	m, _ := camelRuntime.Send("direct:sum", nil, map[string]any{
		"a": 1,
		"b": 39,
	})
	fmt.Printf("%+v\n", m.Payload())

	time.Sleep(1 * time.Minute)

	camelRuntime.Stop()
	println("stop")
}
