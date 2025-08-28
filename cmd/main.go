package main

import (
	"errors"
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

	camelRuntime.RegisterRoute(camel.NewRoute("err", "direct:err",
		processor.DoTry(processor.Process(func(message *camel.Message) {
			panic("zZZZ")
			if !message.HasHeader("a") {
				message.Error = errors.New("Not defined header: a")
			}
		})).Catch(func(err error) bool {
			return err != nil
		}, processor.Process(func(message *camel.Message) {
			fmt.Println(">>>>>", message.Error)
		})),
	))

	camelRuntime.RegisterRoute(camel.NewRoute("sum", "direct:sum",
		processor.Pipeline(
			processor.SetBody(expr.Func(func(message *camel.Message) (any, error) {
				return message.MustHeader("a").(int) + message.MustHeader("b").(int), nil
			})),
			processor.To("direct:print"),
			processor.Choice().
				When(expr.MustSimple("payload == 40"), processor.Process(func(message *camel.Message) {
					message.Body = 444
				})).
				Otherwise(processor.To("direct:x")),
		),
	))

	camelRuntime.RegisterRoute(camel.NewRoute("t", "timer:zzz", processor.Pipeline(
		processor.Process(func(message *camel.Message) {
			message.Body = fmt.Sprintf("COUNT: %d", message.Body)
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
	fmt.Printf("%+v\n", m.Body)

	camelRuntime.Send("direct:err", nil, nil)

	time.Sleep(1 * time.Minute)

	camelRuntime.Stop()
	println("stop")
}
