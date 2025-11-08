package test

import (
	"context"
	"errors"
	"fmt"
	"github.com/paveldanilin/go-camel/pkg/camel"
	"github.com/paveldanilin/go-camel/pkg/camel/component/direct"
	"github.com/paveldanilin/go-camel/pkg/camel/converter"
	"github.com/paveldanilin/go-camel/pkg/camel/env"
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
	"github.com/paveldanilin/go-camel/pkg/camel/expr"
	"testing"
)

func TestRoute_SetBody(t *testing.T) {
	var testCamelRuntime = camel.NewRuntime(camel.RuntimeConfig{Name: "CamelTestRuntime"})
	testCamelRuntime.MustRegisterComponent(direct.NewComponent())

	defer testCamelRuntime.Stop()

	// Build camel step definition
	route, err := camel.NewRoute("sum", "direct:sum").
		Choice("test input data").
		When(expr.Simple("header.a==nil"), func(b *camel.RouteBuilder) {
			b.SetError("", errors.New("not defined mandatory parameter: a"))
		}).
		When(expr.Simple("header.b==nil"), func(b *camel.RouteBuilder) {
			b.SetError("", errors.New("not defined mandatory parameter: b"))
		}).
		Otherwise(func(b *camel.RouteBuilder) {
			b.SetBody("", expr.Simple("header.a + header.b"))
		}).
		Build()
	if err != nil {
		t.Errorf("TestRoute_SetBody(): failed to build 'sum' step: %s", err)
	}

	// Register step in camel runtime
	err = testCamelRuntime.RegisterRoute(route)
	if err != nil {
		t.Errorf("TestRoute_SetBody(): failed to register 'sum' step in runtime: %s", err)
	}

	// Start camel runtime
	err = testCamelRuntime.Start()
	if err != nil {
		t.Errorf("TestRoute_SetBody(): failed to start camel runtime: %s", err)
	}

	tests := []struct {
		name       string
		headers    exchange.Map
		wantErr    bool
		wantResult int
	}{
		{
			name: "Success",
			headers: exchange.Map{
				"a": 2,
				"b": 2,
			},
			wantResult: 4,
			wantErr:    false,
		},
		{
			name: "Not passed mandatory parameter: a",
			headers: exchange.Map{
				"b": 2,
			},
			wantErr: true,
		},
		{
			name: "Not passed mandatory parameter: b",
			headers: exchange.Map{
				"a": 2,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := testCamelRuntime.SendHeaders(context.TODO(), "direct:sum", tt.headers)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestRoute_SetBody() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.wantResult != result.Body {
				t.Errorf("TestRoute_SetBody(): expected result %d, but got %d", tt.wantResult, result.Body)
			}
		})
	}
}

type person struct {
	Name string `json:"name"`
	Age  uint   `json:"age"`
}

func TestRoute_Unmarshal(t *testing.T) {
	var testCamelRuntime = camel.NewRuntime(camel.RuntimeConfig{Name: "CamelTestRuntime"})
	testCamelRuntime.MustRegisterComponent(direct.NewComponent())

	defer testCamelRuntime.Stop()

	// Build camel step definition
	route, err := camel.NewRoute("isPersonAdult", "direct:isPersonAdult").
		Unmarshal("", "json", person{}).
		Choice("test person age").
		// now body contains person{}
		When(expr.Simple("body.Age>=18"), func(b *camel.RouteBuilder) {
			b.SetBody("person is adult", expr.Constant(true))
		}).
		When(expr.Simple("body.Age<18"), func(b *camel.RouteBuilder) {
			b.SetBody("person is not adult", expr.Constant(false))
		}).
		EndChoice().
		Build()
	if err != nil {
		t.Errorf("TestRoute_Unmarshal(): failed to build 'isPersonAdult' step: %s", err)
	}

	// Register step in camel runtime
	err = testCamelRuntime.RegisterRoute(route)
	if err != nil {
		t.Errorf("TestRoute_Unmarshal(): failed to register 'isPersonAdult' step in runtime: %s", err)
	}

	// Start camel runtime
	err = testCamelRuntime.Start()
	if err != nil {
		t.Errorf("TestRoute_Unmarshal(): failed to start camel runtime: %s", err)
	}

	result, err := testCamelRuntime.SendBody(context.TODO(), "direct:isPersonAdult", `{
		"name": "John",
		"age":	22
	}`)
	if err != nil {
		t.Errorf("TestRoute_Unmarshal(): %s", err)
	}

	expectedResult := true
	if expectedResult != result.Body {
		t.Errorf("TestRoute_Unmarshal(): expected result %v, but got %v", expectedResult, result.Body)
	}
}

type operation struct {
	X any
	Y any
}

func TestRoute_ConvertBody(t *testing.T) {
	convReg := converter.NewRegistry()
	convReg.Register(converter.StringToInt())
	convReg.Register(converter.StringToFloat64())
	convReg.Register(converter.Func[map[string]any, *operation](func(v map[string]any, _ map[string]any) (*operation, error) {
		op := &operation{
			X: v["x"],
			Y: v["y"],
		}
		return op, nil
	}))

	var testCamelRuntime = camel.NewRuntime(camel.RuntimeConfig{
		Name:              "CamelTestRuntime",
		ConverterRegistry: convReg,
	})
	testCamelRuntime.MustRegisterComponent(direct.NewComponent())

	defer testCamelRuntime.Stop()

	route, err := camel.NewRoute("convertData", "direct:conv").
		ConvertBodyToNamed("", "*test.operation", nil).
		SetBody("", expr.Simple("body.X + body.Y + 6.2")).
		Build()

	if err != nil {
		t.Errorf("TestRoute_ConvertBody(): failed to build 'convertData' step: %s", err)
	}

	// Register route in camel runtime
	err = testCamelRuntime.RegisterRoute(route)
	if err != nil {
		t.Errorf("TestRoute_ConvertBody(): failed to register 'convertData' step in runtime: %s", err)
	}

	// Start camel runtime
	err = testCamelRuntime.Start()
	if err != nil {
		t.Errorf("TestRoute_ConvertBody(): failed to start camel runtime: %s", err)
	}

	result, err := testCamelRuntime.SendBody(context.TODO(), "direct:conv", map[string]any{
		"x": 1,
		"y": 5,
	})
	if err != nil {
		t.Errorf("TestRoute_ConvertBody(): %s", err)
	}

	fmt.Println(result)
}

func TestRoute_SetHeader(t *testing.T) {
	var testCamelRuntime = camel.NewRuntime(camel.RuntimeConfig{Name: "CamelTestRuntime"})
	testCamelRuntime.MustRegisterComponent(direct.NewComponent())

	defer testCamelRuntime.Stop()

	// Build camel step definition
	route, err := camel.NewRoute("concat", "direct:concat").
		SetHeader("", "year", expr.Constant(2025)).
		SetHeader("", "month", expr.Simple("'11'")).
		SetHeader("", "day", expr.Func(func(_ *exchange.Exchange) (any, error) {
			return 8, nil
		})).
		SetBody("", expr.Simple("toString(header.year) + '.' + header.month + '.' + toString(header.day)")).
		Build()
	if err != nil {
		t.Fatalf("TestRoute_SetHeader(): failed to build 'concat' step: %s", err)
	}

	// Register step in camel runtime
	err = testCamelRuntime.RegisterRoute(route)
	if err != nil {
		t.Fatalf("TestRoute_SetHeader(): failed to register 'concat' step in runtime: %s", err)
	}

	// Start camel runtime
	err = testCamelRuntime.Start()
	if err != nil {
		t.Fatalf("TestRoute_SetHeader(): failed to start camel runtime: %s", err)
	}

	result, err := testCamelRuntime.SendHeaders(context.TODO(), "direct:concat", nil)
	if err != nil {
		t.Fatalf("TestRoute_SetHeader(): failed to call route: %s", err)
	}

	wantResult := "2025.11.8"
	if result.Body != wantResult {
		t.Fatalf("TestRoute_SetHeader(): expected result %v, but got %v", wantResult, result.Body)
	}
}

func TestRoute_Env(t *testing.T) {
	var testCamelRuntime = camel.NewRuntime(camel.RuntimeConfig{
		Name: "CamelTestRuntime",
		Env: env.NewMapEnv(map[string]string{
			"endpointName": "concat",
		}),
	})
	testCamelRuntime.MustRegisterComponent(direct.NewComponent())

	defer testCamelRuntime.Stop()

	// Build camel step definition
	route, err := camel.NewRoute("concat", "direct:${endpointName}").
		SetHeader("", "year", expr.Constant(2025)).
		SetHeader("", "month", expr.Simple("'11'")).
		SetHeader("", "day", expr.Func(func(_ *exchange.Exchange) (any, error) {
			return 8, nil
		})).
		SetBody("", expr.Simple("toString(header.year) + '.' + header.month + '.' + toString(header.day)")).
		Build()
	if err != nil {
		t.Fatalf("TestRoute_SetHeader(): failed to build 'concat' step: %s", err)
	}

	// Register step in camel runtime
	err = testCamelRuntime.RegisterRoute(route)
	if err != nil {
		t.Fatalf("TestRoute_SetHeader(): failed to register 'concat' step in runtime: %s", err)
	}

	// Start camel runtime
	err = testCamelRuntime.Start()
	if err != nil {
		t.Fatalf("TestRoute_SetHeader(): failed to start camel runtime: %s", err)
	}

	result, err := testCamelRuntime.SendHeaders(context.TODO(), "direct:concat", nil)
	if err != nil {
		t.Fatalf("TestRoute_SetHeader(): failed to call route: %s", err)
	}

	wantResult := "2025.11.8"
	if result.Body != wantResult {
		t.Fatalf("TestRoute_SetHeader(): expected result %v, but got %v", wantResult, result.Body)
	}
}
