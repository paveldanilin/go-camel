package test

import (
	"context"
	"errors"
	"github.com/paveldanilin/go-camel/camel"
	"github.com/paveldanilin/go-camel/component/direct"
	"testing"
)

func TestRoute_SetBody(t *testing.T) {
	var testCamelRuntime = camel.NewRuntime(camel.RuntimeConfig{Name: "CamelTestRuntime"})
	testCamelRuntime.MustRegisterComponent(direct.NewComponent())

	defer testCamelRuntime.Stop()

	// Build camel route definition
	route, err := camel.NewRoute("sum", "direct:sum").
		Choice("test input data").
		When(camel.Simple("headers.a==nil"), func(b *camel.RouteBuilder) {
			b.SetError("", errors.New("not defined mandatory parameter: a"))
		}).
		When(camel.Simple("headers.b==nil"), func(b *camel.RouteBuilder) {
			b.SetError("", errors.New("not defined mandatory parameter: b"))
		}).
		Otherwise(func(b *camel.RouteBuilder) {
			b.SetBody("", camel.Simple("headers.a + headers.b"))
		}).
		Build()
	if err != nil {
		t.Errorf("TestRoute_SetBody(): failed to build 'sum' route: %s", err)
	}

	// Register route in camel runtime
	err = testCamelRuntime.RegisterRoute(route)
	if err != nil {
		t.Errorf("TestRoute_SetBody(): failed to register 'sum' route in runtime: %s", err)
	}

	// Start camel runtime
	err = testCamelRuntime.Start()
	if err != nil {
		t.Errorf("TestRoute_SetBody(): failed to start camel runtime: %s", err)
	}

	tests := []struct {
		name       string
		headers    camel.Map
		wantErr    bool
		wantResult int
	}{
		{
			name: "Success",
			headers: camel.Map{
				"a": 2,
				"b": 2,
			},
			wantResult: 4,
			wantErr:    false,
		},
		{
			name: "Not passed mandatory parameter: a",
			headers: camel.Map{
				"b": 2,
			},
			wantErr: true,
		},
		{
			name: "Not passed mandatory parameter: b",
			headers: camel.Map{
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

	// Build camel route definition
	route, err := camel.NewRoute("isPersonAdult", "direct:isPersonAdult").
		//Unmarshal("json", person{}).
		Choice("test person age").
		// now body contains person{}
		When(camel.Simple("body.age>=18"), func(b *camel.RouteBuilder) {
			b.SetBody("person is adult", camel.Constant(true))
		}).
		When(camel.Simple("body.age<18"), func(b *camel.RouteBuilder) {
			b.SetBody("person is not adult", camel.Constant(false))
		}).
		EndChoice().
		Build()
	if err != nil {
		t.Errorf("TestRoute_Unmarshal(): failed to build 'isPersonAdult' route: %s", err)
	}

	// Register route in camel runtime
	err = testCamelRuntime.RegisterRoute(route)
	if err != nil {
		t.Errorf("TestRoute_Unmarshal(): failed to register 'isPersonAdult' route in runtime: %s", err)
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
