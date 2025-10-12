package camel

import (
	"testing"
)

type jsonModelMarshal struct {
	Page   int
	Fruits []string
}

func TestJsonProcessorMarshal(t *testing.T) {
	expectedBody := `{"Page":2,"Fruits":["orange","abricot"]}`
	model := jsonModelMarshal{Page: 2, Fruits: []string{"orange", "abricot"}}
	jsonProcMarshal := newMarshalProcessor("json", model, []byte{})
	exchange := NewExchange(nil, nil)

	jsonProcMarshal.Process(exchange)

	if exchange.Message().Body != expectedBody {
		t.Errorf("TestJsonProcessorMarshal() = %v; want body %v", exchange.Message().Body, expectedBody)
	}
}
