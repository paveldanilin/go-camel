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
		processor.SetBody(expr.MustSimple("headers.a + headers.b")),
	))

	// Ticks every 5 seconds
	camelRuntime.MustRegisterRoute(camel.NewRoute("t", "timer:myTimer?interval=5s", processor.PipelineWithConfig(
		processor.PipelineConfig{
			FailFast: false, // keep processing pipeline in case of any error
		},
		processor.SetBody(expr.MustSimple("'COUNT: ' + string(headers.CamelTimerCounter)")),
		processor.LogMessage(">>"),
	)))

	// Start Camel runtime
	err := camelRuntime.Start()
	if err != nil {
		panic(err)
	}

	// Calc sum
	m, _ := camelRuntime.Send("direct:sum", nil, map[string]any{
		"a": 1,
		"b": 39,
	})
	fmt.Printf("SUM: %+v\n", m.Body)

	time.Sleep(1 * time.Minute)

	camelRuntime.Stop()
	println("stop")
}
