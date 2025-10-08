package main

import (
	"context"
	"fmt"
	"github.com/paveldanilin/go-camel/camel"
	"github.com/paveldanilin/go-camel/component/direct"
	"github.com/paveldanilin/go-camel/component/timer"
	"strings"
	"time"
)

func main() {

	camelRuntime := camel.NewRuntime(camel.RuntimeConfig{
		ExchangeFactory: nil,
		FuncRegistry:    nil,
	})
	camelRuntime.MustRegisterComponent(direct.NewComponent())
	camelRuntime.MustRegisterComponent(timer.NewComponent())

	camelRuntime.MustRegisterFunc("x10Func", func(exchange *camel.Exchange) {
		exchange.Message().Body = exchange.Message().Body.(int) * 100
	})

	r, err := camel.NewRoute("sum", "direct:sum").
		SetBody("calc sum", camel.Simple("headers.a + headers.b")).
		Choice("test sum result").
		When(camel.Simple("body == 40"), func(b *camel.RouteBuilder) {
			b.Sleep("", 2500)
			b.SetBody("double body value", camel.Simple("body * 2"))
		}).
		Otherwise(func(b *camel.RouteBuilder) {
			b.SetBody("x*4", camel.Simple("body * 4"))
		}).
		Try("", func(b *camel.RouteBuilder) {
			b.Func("x10", "x10Func")
		}).
		Catch(camel.ErrAny(), func(b *camel.RouteBuilder) {
			b.SetBody("xxx", camel.Simple("'>>' + error.Error() + '<<'"))
		}).
		EndTry().
		Multicast("multi tasks").ParallelProcessing().
		Output(func(b *camel.RouteBuilder) {
			b.Sleep("", 15000)
			b.LogWarn("", "xxx> ${body}")
		}).
		Output(func(b *camel.RouteBuilder) {
			b.Sleep("", 5000)
			b.SetProperty("setGame", "game", camel.Constant("DooM"))
			b.LogInfo("", "yy>>${body}>>${properties.game}")
		}).
		EndMulticast().
		Build()
	if err != nil {
		panic(err)
	}

	/*getDepth := func(step camel.RouteStep, depth int) error {
		fmt.Printf("%s> %s [%T]\n", strings.Repeat("-", depth+1), step.StepName(), step)
		return nil
	}
	_ = camel.WalkRoute(r, getDepth)*/

	camelRuntime.MustRegisterRoute(r)

	// Ticks every 5 seconds
	/*
		r, err = camel.NewRoute("ticker", "timer:myTimer?interval=5s").
			Pipeline("on tick", true, func(b *dsl.RouteBuilder) {
				b.SetBody("set count", camel.Simple("'COUNT: ' + string(headers.CamelTimerCounter)"))
			}).
			Build()
		if err != nil {
			panic(err)
		}
		camelRuntime.MustRegisterRoute(r)
	*/

	// Start Camel runtime
	err = camelRuntime.Start()
	if err != nil {
		panic(err)
	}

	// Calc sum
	ctx, _ := context.WithTimeout(context.Background(), 60*time.Second)
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
