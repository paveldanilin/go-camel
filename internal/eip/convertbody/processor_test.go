package convertbody

import (
	"errors"
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
	"reflect"
	"strconv"
	"testing"
)

type string2int struct {
}

func (c *string2int) Convert(v any, _ reflect.Type, _ map[string]any) (any, error) {
	switch t := v.(type) {
	case string:
		return strconv.Atoi(t)
	}
	return v, errors.New("cannot convert value")
}

func TestConvertBodyProcessor(t *testing.T) {
	p := NewProcessor("test", "test", reflect.TypeOf(0), nil, &string2int{})
	e := exchange.NewExchange(nil)
	e.Message().Body = "333"

	p.Process(e)

	expectedValue := 333
	if e.Message().Body != expectedValue {
		t.Fatalf("TestConvertBodyProcessor() = %v; want %d", e.Message().Body, expectedValue)
	}
}
