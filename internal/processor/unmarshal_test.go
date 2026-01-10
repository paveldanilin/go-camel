package processor

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
	p := NewUnmarshal(jsonModelUnmarshal{}, dataformat.JSON{})

	e := exchange.NewExchange(nil)
	e.Message().SetBody(`{"Page": 1, "Fruits": ["apple", "peach"]}`)

	p.Process(e)

	expectedBody := reflect.TypeOf(&jsonModelUnmarshal{})
	if reflect.TypeOf(e.Message().Body()) != expectedBody {
		t.Errorf("TestUnmarshalProcessorJson() = %v; want body %v", e.Message().Body(), expectedBody)
	}
}
