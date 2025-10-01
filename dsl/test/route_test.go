package test

import (
	"fmt"
	"github.com/paveldanilin/go-camel/dsl"
	"strings"
	"testing"
)

func TestRouteBuilder(t *testing.T) {
	r, err := dsl.NewRouteBuilder("sum", "direct:sum").
		SetHeader("set a", "a", dsl.Constant(1)).
		SetHeader("set b", "b", dsl.Constant(1)).
		SetBody("set body sum", dsl.Simple("header.a + header.b")).
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
	r, err := dsl.NewRouteBuilder("test age", "direct:age").
		SetHeader("set age", "age", dsl.Constant(10)).
		Choice("test age").
		When(dsl.Simple("header.age < 14"), func(b *dsl.RouteBuilder) {
			b.SetBody("set access", dsl.Constant("DENY"))
		}).
		When(dsl.Simple("header.age >= 15"), func(b *dsl.RouteBuilder) {
			b.SetBody("set access", dsl.Constant("ALLOW"))
		}).
		EndChoice().
		SetHeader("set access", "access", dsl.Simple("message.body")).
		Choice("test access").
		When(dsl.Simple("header.access == 'ALLOW'"), func(b *dsl.RouteBuilder) {
			b.SetBody("set data link", dsl.Constant("http://secret.data.link"))
		}).
		Otherwise(func(b *dsl.RouteBuilder) {
			b.SetBody("set forbidden message", dsl.Constant("Access denied, bye"))
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
	r, err := dsl.NewRouteBuilder("test deep nested dsl", "direct:deep-nested").
		Pipeline("pipeline_1", false, func(b *dsl.RouteBuilder) {
			b.Pipeline("pipeline_2", false, func(b *dsl.RouteBuilder) {
				b.Pipeline("pipeline_3", false, func(b *dsl.RouteBuilder) {
					b.Choice("choice_1").
						When(dsl.Simple("true"), func(b *dsl.RouteBuilder) {
							b.SetHeader("", "xxx", dsl.Constant("yyy"))
						}).
						When(dsl.Simple("1==1"), func(b *dsl.RouteBuilder) {
							b.Choice("choice_2").
								When(dsl.Simple("2==2"), func(b *dsl.RouteBuilder) {
									b.Choice("choice_3").
										When(dsl.Simple("3==3"), func(b *dsl.RouteBuilder) {
											b.SetBody("set boy", dsl.Constant(1))
											b.Pipeline("props", false, func(b *dsl.RouteBuilder) {
												b.SetProperty("", "x", dsl.Constant("x"))
												b.SetProperty("", "y", dsl.Constant("y"))
												b.SetProperty("", "z", dsl.Constant("z"))
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
	getDepth := func(step dsl.RouteStep, depth int) error {
		routeDepth = depth
		fmt.Printf("%s> %s [%T]\n", strings.Repeat("-", depth+1), step.StepName(), step)
		return nil
	}
	_ = dsl.WalkRoute(r, getDepth)

	if routeDepth != 10 {
		t.Errorf("expected depth of dsl is 10, but got %d", routeDepth)
	}
}

func TestRouteBuilder_DoTry(t *testing.T) {
	r, err := dsl.NewRouteBuilder("doTry", "direct:doTry").
		SetBody("set empty body", dsl.Constant("")).
		Try("safe block", func(b *dsl.RouteBuilder) {
			b.To("critical operation", "http://api.secret.com?key=xyz&httpMethod=GET")
		}).
		Catch(dsl.ErrEquals("io errors"), func(b *dsl.RouteBuilder) {
			b.SetProperty("error", "io.error", dsl.Constant("IO error"))
		}).
		Catch(dsl.ErrEquals("net error"), func(b *dsl.RouteBuilder) {
			b.SetProperty("error", "net.error", dsl.Constant("NET error"))
		}).
		EndTry().
		Choice("if error").
		When(dsl.Simple("property.error != nil"), func(b *dsl.RouteBuilder) {
			b.Try("safe send error", func(b *dsl.RouteBuilder) {
				b.SetBody("set error body", dsl.Simple("property.error"))
				b.To("send error to collector", "http://error.collector?httpMethod=POST")
			}).Catch(dsl.ErrAny(), func(b *dsl.RouteBuilder) {
				b.SetHeader("set error", "error", dsl.Simple("message.body"))
			}).Finally(func(b *dsl.RouteBuilder) {
				b.SetBody("set finally body", dsl.Constant("FIN"))
			})
		}).
		EndChoice().
		Build()

	if err != nil {
		t.Error(err)
	}

	dsl.WalkRoute(r, func(step dsl.RouteStep, depth int) error {
		fmt.Printf("[%s]\n", step.StepName())
		return nil
	})
}

func TestRouteBuilder_LoopWhile(t *testing.T) {
	r, err := dsl.NewRouteBuilder("loop", "direct:loopWhile").
		SetBody("set data", dsl.Simple("[1,2,3,4,5]")).
		LoopWhile("find digit", dsl.Simple("CAMEL_LOOP_INDEX < 10"), false, func(b *dsl.RouteBuilder) {
			b.SetBody("set decorated body", dsl.Simple("[properties.CAMEL_LOOP_INDEX]"))
			b.Choice("test digit").
				When(dsl.Simple("message.body==3"), func(b *dsl.RouteBuilder) {
					b.SetProperty("exit loop", "CAMEL_LOOP_BREAK", dsl.Constant(true))
				}).
				EndChoice()
		}).
		Build()

	if err != nil {
		t.Error(err)
	}

	dsl.WalkRoute(r, func(step dsl.RouteStep, depth int) error {
		fmt.Printf("[%s]\n", step.StepName())
		return nil
	})
}
