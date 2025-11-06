package unmarshal

import (
	"github.com/paveldanilin/go-camel/pkg/camel/dataformat"
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
	"reflect"
	"testing"
)

type jsonModelUnmarshal struct {
	Page   int
	Fruits []string
}

func TestUnmarshalProcessorJson(t *testing.T) {
	p := NewProcessor("", "", jsonModelUnmarshal{}, dataformat.JSON{})

	e := exchange.NewExchange(nil)
	e.Message().Body = `{"Page": 1, "Fruits": ["apple", "peach"]}`

	p.Process(e)

	expectedBody := reflect.TypeOf(&jsonModelUnmarshal{})
	if reflect.TypeOf(e.Message().Body) != expectedBody {
		t.Errorf("TestJsonProcessorUnmarshal() = %v; want body %v", e.Message().Body, expectedBody)
	}
}
