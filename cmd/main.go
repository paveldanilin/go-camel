package main

import (
	"fmt"
	"github.com/paveldanilin/go-camel/camel"
	"github.com/paveldanilin/go-camel/camel/component/direct"
	"github.com/paveldanilin/go-camel/camel/component/timer"
	"github.com/paveldanilin/go-camel/camel/expr"
	"github.com/paveldanilin/go-camel/camel/processor"
	"log"
	"time"
)

func main() {

	camelRuntime := camel.NewRuntime()
	camelRuntime.MustRegisterComponent(direct.NewComponent())
	camelRuntime.MustRegisterComponent(timer.NewComponent())

	camelRuntime.MustRegisterRoute(camel.NewRoute("err", "direct:err",
		processor.DoTry(processor.Process(func(exchange *camel.Exchange) {
			panic("zZZZ")
		})).Catch(camel.ErrorAny(), processor.Process(func(exchange *camel.Exchange) {
			fmt.Println(">>>>>", exchange.Error)
		})),
	))

	camelRuntime.MustRegisterRoute(camel.NewRoute("sum", "direct:sum",
		processor.Pipeline(
			processor.SetBody(expr.Func(func(exchange *camel.Exchange) (any, error) {
				return exchange.Message().MustHeader("a").(int) + exchange.Message().MustHeader("b").(int), nil
			})),
			processor.To("direct:print"),
			processor.Choice().
				When(expr.MustSimple("body == 40"), processor.Process(func(exchange *camel.Exchange) {
					exchange.Message().Body = 444
				})).
				Otherwise(processor.To("direct:x")),
		),
	))

	camelRuntime.MustRegisterRoute(camel.NewRoute("t", "timer:myTimer?interval=5s", processor.PipelineWithConfig(
		processor.PipelineConfig{
			FailFast: false, // keep processing pipeline in case of any error
		},
		processor.Process(func(exchange *camel.Exchange) {
			exchange.Message().Body = fmt.Sprintf("COUNT: %d", exchange.Message().Body)
		}),
		processor.To("direct:gg"),
		processor.Process(func(exchange *camel.Exchange) {
			if exchange.Error != nil {
				log.Printf("error: %s", exchange.Error)
			}
		}),
	)))

	camelRuntime.MustRegisterRoute(camel.NewRoute("print", "direct:print", processor.LogMessage("print-1")))

	camelRuntime.MustRegisterRoute(camel.NewRoute("print1", "direct:print", processor.LogMessage("print-2")))

	camelRuntime.MustRegisterRoute(camel.NewRoute("x", "direct:x", processor.LogMessage("xxx")))

	err := camelRuntime.Start()
	if err != nil {
		panic(err)
	}

	//m, _ := camelRuntime.Send("direct:sum", nil, map[string]any{
	//	"a": 1,
	//	"b": 39,
	//})
	//fmt.Printf("%+v\n", m.Body)

	//_, _ = camelRuntime.Send("direct:err", nil, nil)

	time.Sleep(1 * time.Minute)

	camelRuntime.Stop()
	println("stop")
}
