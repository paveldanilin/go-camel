package camel

import (
	"reflect"
	"testing"
)

type jsonModel struct {
	Page   int
	Fruits []string
}

func TestJsonProcessorMarshal(t *testing.T) {
	expectedBody := `{"Page":2,"Fruits":["orange","abricot"]}`
	model := jsonModel{Page: 2, Fruits: []string{"orange", "abricot"}}
	jsonProcMarshal := newJsonProcessor(Marshal, model, []byte{})
	exchange := NewExchange(nil, nil)

	jsonProcMarshal.Process(exchange)

	if exchange.Message().Body != expectedBody {
		t.Errorf("TestJsonProcessorMarshal() = %v; want body %v", exchange.Message().Body, expectedBody)
	}
}

func TestJsonProcessorUnmarshal(t *testing.T) {
	expectedBody := reflect.TypeOf(jsonModel{})
	json := `{"Page": 1, "Fruits": ["apple", "peach"]}`
	jsonProcUnmarshal := newJsonProcessor(Unmarshal, jsonModel{}, []byte(json))
	exchange := NewExchange(nil, nil)

	jsonProcUnmarshal.Process(exchange)

	if reflect.TypeOf(exchange.Message().Body) != expectedBody {
		t.Errorf("TestJsonProcessorUnmarshal() = %v; want body %v", exchange.Message().Body, expectedBody)
	}
}
