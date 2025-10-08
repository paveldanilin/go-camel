package test

import (
	"fmt"
	"github.com/paveldanilin/go-camel/camel"
	"strings"
	"testing"
)

func TestRouteBuilder(t *testing.T) {
	r, err := camel.NewRoute("sum", "direct:sum").
		SetHeader("set a", "a", camel.Constant(1)).
		SetHeader("set b", "b", camel.Constant(1)).
		SetBody("set body sum", camel.Simple("header.a + header.b")).
		Build()

	if err != nil {
		t.Error(err)
	}

	if r.Name != "sum" {
		t.Errorf("expected dsl name 'sum', but got '%s'", r.Name)
	}

	if r.From != "direct:sum" {
		t.Errorf("expected dsl from 'direct:sum', but got '%s'", r.From)
	}

	if len(r.Steps) != 3 {
		t.Errorf("expected 3 steps, but got %d", len(r.Steps))
	}

	if r.Steps[0].StepName() != "set a" {
		t.Errorf("expected first step 'set a', but got '%s'", r.Steps[0].StepName())
	}

	if r.Steps[1].StepName() != "set b" {
		t.Errorf("expected first step 'set b', but got '%s'", r.Steps[1].StepName())
	}

	if r.Steps[2].StepName() != "set body sum" {
		t.Errorf("expected first step 'set body sum', but got '%s'", r.Steps[2].StepName())
	}
}

func TestRouteBuilder_Choice(t *testing.T) {
	r, err := camel.NewRoute("test age", "direct:age").
		SetHeader("set age", "age", camel.Constant(10)).
		Choice("test age").
		When(camel.Simple("header.age < 14"), func(b *camel.RouteBuilder) {
			b.SetBody("set access", camel.Constant("DENY"))
		}).
		When(camel.Simple("header.age >= 15"), func(b *camel.RouteBuilder) {
			b.SetBody("set access", camel.Constant("ALLOW"))
		}).
		EndChoice().
		SetHeader("set access", "access", camel.Simple("message.body")).
		Choice("test access").
		When(camel.Simple("header.access == 'ALLOW'"), func(b *camel.RouteBuilder) {
			b.SetBody("set data link", camel.Constant("http://secret.data.link"))
		}).
		Otherwise(func(b *camel.RouteBuilder) {
			b.SetBody("set forbidden message", camel.Constant("Access denied, bye"))
		}).
		Build()

	if err != nil {
		t.Error(err)
	}

	if len(r.Steps) != 4 {
		t.Errorf("expected 4 steps, but got %d", len(r.Steps))
	}
}

func TestRouteBuilder_DeepNested(t *testing.T) {
	r, err := camel.NewRoute("test deep nested dsl", "direct:deep-nested").
		Pipeline("pipeline_1", false, func(b *camel.RouteBuilder) {
			b.Pipeline("pipeline_2", false, func(b *camel.RouteBuilder) {
				b.Pipeline("pipeline_3", false, func(b *camel.RouteBuilder) {
					b.Choice("choice_1").
						When(camel.Simple("true"), func(b *camel.RouteBuilder) {
							b.SetHeader("", "xxx", camel.Constant("yyy"))
						}).
						When(camel.Simple("1==1"), func(b *camel.RouteBuilder) {
							b.Choice("choice_2").
								When(camel.Simple("2==2"), func(b *camel.RouteBuilder) {
									b.Choice("choice_3").
										When(camel.Simple("3==3"), func(b *camel.RouteBuilder) {
											b.SetBody("set boy", camel.Constant(1))
											b.Pipeline("props", false, func(b *camel.RouteBuilder) {
												b.SetProperty("", "x", camel.Constant("x"))
												b.SetProperty("", "y", camel.Constant("y"))
												b.SetProperty("", "z", camel.Constant("z"))
											})
										})
								})

						})
				})
			})
		}).Build()

	if err != nil {
		t.Error(err)
	}

	var routeDepth = 0
	getDepth := func(step camel.RouteStep, depth int) error {
		routeDepth = depth
		fmt.Printf("%s> %s [%T]\n", strings.Repeat("-", depth+1), step.StepName(), step)
		return nil
	}
	_ = camel.WalkRoute(r, getDepth)

	if routeDepth != 10 {
		t.Errorf("expected depth of dsl is 10, but got %d", routeDepth)
	}
}

func TestRouteBuilder_DoTry(t *testing.T) {
	r, err := camel.NewRoute("doTry", "direct:doTry").
		SetBody("set empty body", camel.Constant("")).
		Try("safe block", func(b *camel.RouteBuilder) {
			b.To("critical operation", "http://api.secret.com?key=xyz&httpMethod=GET")
		}).
		Catch(camel.ErrEquals("io errors"), func(b *camel.RouteBuilder) {
			b.SetProperty("error", "io.error", camel.Constant("IO error"))
		}).
		Catch(camel.ErrEquals("net error"), func(b *camel.RouteBuilder) {
			b.SetProperty("error", "net.error", camel.Constant("NET error"))
		}).
		EndTry().
		Choice("if error").
		When(camel.Simple("property.error != nil"), func(b *camel.RouteBuilder) {
			b.Try("safe send error", func(b *camel.RouteBuilder) {
				b.SetBody("set error body", camel.Simple("property.error"))
				b.To("send error to collector", "http://error.collector?httpMethod=POST")
			}).Catch(camel.ErrAny(), func(b *camel.RouteBuilder) {
				b.SetHeader("set error", "error", camel.Simple("message.body"))
			}).Finally(func(b *camel.RouteBuilder) {
				b.SetBody("set finally body", camel.Constant("FIN"))
			})
		}).
		EndChoice().
		Build()

	if err != nil {
		t.Error(err)
	}

	camel.WalkRoute(r, func(step camel.RouteStep, depth int) error {
		fmt.Printf("[%s]\n", step.StepName())
		return nil
	})
}

func TestRouteBuilder_LoopWhile(t *testing.T) {
	r, err := camel.NewRoute("loop", "direct:loopWhile").
		SetBody("set data", camel.Simple("[1,2,3,4,5]")).
		LoopWhile("find digit", camel.Simple("CAMEL_LOOP_INDEX < 10"), false, func(b *camel.RouteBuilder) {
			b.SetBody("set decorated body", camel.Simple("[properties.CAMEL_LOOP_INDEX]"))
			b.Choice("test digit").
				When(camel.Simple("message.body==3"), func(b *camel.RouteBuilder) {
					b.SetProperty("exit loop", "CAMEL_LOOP_BREAK", camel.Constant(true))
				}).
				EndChoice()
		}).
		Build()

	if err != nil {
		t.Error(err)
	}

	camel.WalkRoute(r, func(step camel.RouteStep, depth int) error {
		fmt.Printf("[%s]\n", step.StepName())
		return nil
	})
}
