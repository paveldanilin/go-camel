package main

import (
	"context"
	"fmt"
	"github.com/paveldanilin/go-camel/pkg/camel"
	"github.com/paveldanilin/go-camel/pkg/camel/component/direct"
	"github.com/paveldanilin/go-camel/pkg/camel/component/timer"
	"github.com/paveldanilin/go-camel/pkg/camel/errs"
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
	"github.com/paveldanilin/go-camel/pkg/camel/expr"
	"time"
)

func main() {

	camelRuntime := camel.NewRuntime(camel.RuntimeConfig{
		ExchangeFactory: nil,
		FuncRegistry:    nil,
	})
	camelRuntime.MustRegisterComponent(direct.NewComponent())
	camelRuntime.MustRegisterComponent(timer.NewComponent())

	camelRuntime.MustRegisterFunc("x10Func", func(e *exchange.Exchange) {
		e.Message().Body = e.Message().Body.(int) * 100
	})

	r, err := camel.NewRoute("sum", "direct:sum").
		SetBody("calc sum", expr.Simple("header.a + header.b")).
		Choice("test sum result").
		When(expr.Simple("body == 40"), func(b *camel.RouteBuilder) {
			b.Delay("", 2500)
			b.SetBody("double body value", expr.Simple("body * 2"))
		}).
		Otherwise(func(b *camel.RouteBuilder) {
			b.SetBody("x*4", expr.Simple("body * 4"))
		}).
		Try("", func(b *camel.RouteBuilder) {
			b.Func("x10", "x10Func")
		}).
		Catch(errs.Any(), func(b *camel.RouteBuilder) {
			b.SetBody("xxx", expr.Simple("'>>' + error.Error() + '<<'"))
		}).
		EndTry().
		Multicast("multi tasks").ParallelProcessing().
		Process(func(b *camel.RouteBuilder) {
			b.Delay("", 15000)
			b.LogWarn("", "xxx> ${body}")
		}).
		Process(func(b *camel.RouteBuilder) {
			b.Delay("", 5000)
			b.SetProperty("setGame", "game", expr.Constant("DooM"))
			b.LogInfo("", "yy>>${body}>>${property.game}")
		}).
		EndMulticast().
		Build()
	if err != nil {
		panic(err)
	}

	/*getDepth := fn(step camel.RouteStep, depth int) error {
		fmt.Printf("%s> %s [%T]\n", strings.Repeat("-", depth+1), step.StepName(), step)
		return nil
	}
	_ = camel.WalkRoute(r, getDepth)*/

	camelRuntime.MustRegisterRoute(r)

	// Ticks every 5 seconds
	/*
		r, errs = camel.NewRoute("ticker", "timer:myTimer?interval=5s").
			Pipeline("on tick", true, fn(b *dsl.RouteBuilder) {
				b.SetBody("set count", camel.Simple("'COUNT: ' + string(headers.CamelTimerCounter)"))
			}).
			Build()
		if errs != nil {
			panic(errs)
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
	ex, err := camelRuntime.Send(ctx, "direct:sum", 0, exchange.Map{
		"a": 1,
		"b": 39,
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("RESULT: %v+", ex.Message().Body)
	//fmt.Printf("%s > SUM: %+v [%T]\n", strings.Join(ex.MessagePath(), "/"), ex.Message().Body, ex.Message().Body)

	time.Sleep(1 * time.Minute)

	camelRuntime.Stop()
	println("stop")
}
