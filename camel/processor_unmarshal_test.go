package camel

import (
	"github.com/paveldanilin/go-camel/dataformat"
	"reflect"
	"testing"
)

type jsonModelUnmarshal struct {
	Page   int
	Fruits []string
}

func TestUnmarshalProcessorJson(t *testing.T) {
	p := newUnmarshalProcessor("", jsonModelUnmarshal{}, dataformat.JSONFormat{})

	exchange := NewExchange(nil, nil)
	exchange.Message().Body = `{"Page": 1, "Fruits": ["apple", "peach"]}`

	p.Process(exchange)

	expectedBody := reflect.TypeOf(&jsonModelUnmarshal{})
	if reflect.TypeOf(exchange.Message().Body) != expectedBody {
		t.Errorf("TestJsonProcessorUnmarshal() = %v; want body %v", exchange.Message().Body, expectedBody)
	}
}
