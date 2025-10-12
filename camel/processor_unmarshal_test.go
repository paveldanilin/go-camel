package camel

import (
	"reflect"
	"testing"
)

type jsonModelUnmarshal struct {
	Page   int
	Fruits []string
}

func TestUnmarshalProcessorJson(t *testing.T) {
	expectedBody := reflect.TypeOf(jsonModelUnmarshal{})
	json := `{"Page": 1, "Fruits": ["apple", "peach"]}`
	jsonProcUnmarshal := newUnmarshalProcessor("json", jsonModelUnmarshal{}, []byte(json))
	exchange := NewExchange(nil, nil)

	jsonProcUnmarshal.Process(exchange)

	if reflect.TypeOf(exchange.Message().Body) != expectedBody {
		t.Errorf("TestJsonProcessorUnmarshal() = %v; want body %v", exchange.Message().Body, expectedBody)
	}
}
