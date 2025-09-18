package main

import (
	"context"
	"fmt"
	"github.com/paveldanilin/go-camel/camel"
	"github.com/paveldanilin/go-camel/camel/component/direct"
	"github.com/paveldanilin/go-camel/camel/component/timer"
	"github.com/paveldanilin/go-camel/camel/dsl"
	"time"
)

func main() {

	camelRuntime := camel.NewRuntime()
	camelRuntime.MustRegisterComponent(direct.NewComponent())
	camelRuntime.MustRegisterComponent(timer.NewComponent())

	r, err := camel.NewRoute("sum", "direct:sum").
		SetBody("calc sum", camel.Simple("headers.a + headers.b")).
		Build()
	if err != nil {
		panic(err)
	}

	camelRuntime.MustRegisterRoute(r)

	// Ticks every 5 seconds
	r, err = camel.NewRoute("ticker", "timer:myTimer?interval=5s").
		Pipeline("on tick", func(b *dsl.RouteBuilder) {
			b.SetBody("set count", camel.Simple("'COUNT: ' + string(headers.CamelTimerCounter)"))
		}).
		Build()
	if err != nil {
		panic(err)
	}
	camelRuntime.MustRegisterRoute(r)

	// Start Camel runtime
	err = camelRuntime.Start()
	if err != nil {
		panic(err)
	}

	// Calc sum
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	m, err := camelRuntime.SendHeaders(ctx, "direct:sum", camel.Map{
		"a": 1,
		"b": 39,
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("SUM: %+v\n", m.Body)

	time.Sleep(1 * time.Minute)

	camelRuntime.Stop()
	println("stop")
}
