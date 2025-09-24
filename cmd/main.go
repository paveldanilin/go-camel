package main

import (
	"context"
	"fmt"
	"github.com/paveldanilin/go-camel/camel"
	"github.com/paveldanilin/go-camel/component/direct"
	"github.com/paveldanilin/go-camel/component/timer"
	"github.com/paveldanilin/go-camel/dsl"
	"strings"
	"time"
)

func main() {

	camelRuntime := camel.NewRuntime()
	camelRuntime.MustRegisterComponent(direct.NewComponent())
	camelRuntime.MustRegisterComponent(timer.NewComponent())

	camelRuntime.RegisterFunc("x10Func", func(exchange *camel.Exchange) {
		exchange.Message().Body = exchange.Message().Body.(int) * 100
	})

	r, err := camel.NewRoute("sum", "direct:sum").
		SetBody("calc sum", camel.Simple("headers.a + headers.b")).
		Choice("test sum result").
		When(camel.Simple("body == 40"), func(b *dsl.RouteBuilder) {
			b.Sleep("", 2500)
			b.SetBody("double body value", camel.Simple("body * 2"))
		}).
		Otherwise(func(b *dsl.RouteBuilder) {
			b.SetBody("x*4", camel.Simple("body * 4"))
		}).
		Try("", func(b *dsl.RouteBuilder) {
			// Inline function
			//b.Func("x10", func(exchange *camel.Exchange) {
			//	exchange.Message().Body = exchange.Message().Body.(int) * 100
			//})

			// Stored function (register it first RegistrFunc)
			b.Func("x10", "x10Func")

			//b.Func("panic", func(exchange *camel.Exchange) {
			//	panic(errors.New("!panic!"))
			//})
			//b.Func("zzzz", func(exchange *camel.Exchange) {
			//	println("\n!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!\n")
			//})
		}).
		Catch(camel.ErrAny(), func(b *dsl.RouteBuilder) {
			b.SetBody("xxx", camel.Simple("'>>' + error.Error() + '<<'"))
		}).
		EndTry().
		Build()
	if err != nil {
		panic(err)
	}

	camelRuntime.MustRegisterRoute(r)

	// Ticks every 5 seconds
	r, err = camel.NewRoute("ticker", "timer:myTimer?interval=5s").
		Pipeline("on tick", true, func(b *dsl.RouteBuilder) {
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
	ex, err := camelRuntime.Send(ctx, "direct:sum", 0, camel.Map{
		"a": 1,
		"b": 39,
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s > SUM: %+v [%T]\n", strings.Join(ex.MessagePath(), "/"), ex.Message().Body, ex.Message().Body)

	time.Sleep(1 * time.Minute)

	camelRuntime.Stop()
	println("stop")
}
