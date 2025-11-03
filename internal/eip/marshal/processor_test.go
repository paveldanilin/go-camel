package marshal

import (
	"github.com/paveldanilin/go-camel/pkg/camel/dataformat"
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
	"testing"
)

type jsonModelMarshal struct {
	Page   int
	Fruits []string
}

func TestJsonProcessorMarshal(t *testing.T) {
	p := NewProcessor("", "", &dataformat.JSON2{})

	e := exchange.NewExchange(nil)
	e.Message().Body = jsonModelMarshal{Page: 2, Fruits: []string{"orange", "abricot"}}

	p.Process(e)

	if e.IsError() {
		t.Fatal(e.Error())
	}

	expectedBody := `{"Page":2,"Fruits":["orange","abricot"]}`
	if e.Message().Body != expectedBody {
		t.Errorf("TestJsonProcessorMarshal() = [%v]; want body %v", e.Message().Body, expectedBody)
	}
}
