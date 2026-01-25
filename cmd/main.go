package main

import (
	"context"
	"fmt"
	"github.com/paveldanilin/go-camel/pkg/camel"
	"github.com/paveldanilin/go-camel/pkg/camel/component/direct"
	"github.com/paveldanilin/go-camel/pkg/camel/component/timer"
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
	"github.com/paveldanilin/go-camel/pkg/camel/route"
	"github.com/paveldanilin/go-camel/pkg/camel/route/step"
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
		e.Message().SetBody(e.Message().Body().(int) * 100)
	})

	r, err := camel.NewRoute("sum", "direct:sum").
		SetBodyExpr("header.a + header.b", step.WithName("CalcSum")).
		Choice().
		WhenExpr("body == 40", func(b *route.Builder) {
			b.Delay(2500, step.WithName("WaitForTimeout"))
			b.SetBodyExpr("body * 2")
		}).
		Otherwise(route.SetBodyExpr("body * 4")).
		Try(route.Call("x10Func")).
		CatchAny(route.SetBodyExpr("'>>' + error.Error() + '<<'")).
		EndTry().
		Multicast().WithParallelProcessing().
		Process(func(b *route.Builder) {
			b.Delay(15000)
			b.LogInfo("xxx> ${body}")
		}).
		Process(func(b *route.Builder) {
			b.Delay(5000)
			// TODO: parallel multicast sync exchange?
			b.SetPropertyConstant("game", "DooM")
			b.LogInfo("yy>>${body}>>${property.game}")
		}).
		EndMulticast().
		//Func(func(e *exchange.Exchange) {
		//	panic("011001101")
		//}, step.WithName("TestPanic")).
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
		r, errmatcher = camel.NewRoute("ticker", "timer:myTimer?interval=5s").
			Pipeline("on tick", true, fn(b *dsl.RouteBuilder) {
				b.SetBody("set count", camel.Simple("'COUNT: ' + string(headers.CamelTimerCounter)"))
			}).
			Build()
		if errmatcher != nil {
			panic(errmatcher)
		}
		camelRuntime.MustRegisterRoute(r)
	*/
	r, _ = camel.NewRoute("x", "direct:x").
		Convert().Body().ToType("int").EndConvert().
		Build()
	camelRuntime.MustRegisterRoute(r)

	// Start Camel runtime
	err = camelRuntime.Start()
	if err != nil {
		panic(err)
	}

	// ---
	res, _ := camelRuntime.Send(context.TODO(), "direct:x", "123", nil)
	fmt.Printf(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> %v [%T]", res.MessageBody(), res.MessageBody())
	// ---

	// Calc sum
	ctx, _ := context.WithTimeout(context.Background(), 60*time.Second)
	ex, err := camelRuntime.Send(ctx, "direct:sum", 0, exchange.MessageHeaders{
		"a": 1,
		"b": 39,
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("RESULT: %v+", ex.Message().Body())
	//fmt.Printf("%s > SUM: %+v [%T]\n", strings.Join(ex.MessagePath(), "/"), ex.Message().Body, ex.Message().Body)

	time.Sleep(1 * time.Minute)

	camelRuntime.Stop()
	println("stop")
}
