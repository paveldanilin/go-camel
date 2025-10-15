package camel

import (
	"github.com/paveldanilin/go-camel/dataformat"
	"testing"
)

type jsonModelMarshal struct {
	Page   int
	Fruits []string
}

func TestJsonProcessorMarshal(t *testing.T) {
	p := newMarshalProcessor("", &dataformat.JSONFormat{})

	exchange := NewExchange(nil, nil)
	exchange.Message().Body = jsonModelMarshal{Page: 2, Fruits: []string{"orange", "abricot"}}

	p.Process(exchange)

	if exchange.err != nil {
		t.Fatal(exchange.err)
	}

	expectedBody := `{"Page":2,"Fruits":["orange","abricot"]}`
	if exchange.Message().Body != expectedBody {
		t.Errorf("TestJsonProcessorMarshal() = [%v]; want body %v", exchange.Message().Body, expectedBody)
	}
}
