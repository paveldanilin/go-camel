package processor

import (
	"context"
	"errors"
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
	"github.com/stretchr/testify/assert"
	"reflect"
	"strconv"
	"testing"
)

type string2int struct{}

func (c *string2int) Convert(v any, _ reflect.Type, _ map[string]any) (any, error) {
	switch t := v.(type) {
	case string:
		return strconv.Atoi(t)
	}
	return v, errors.New("cannot convert value")
}

func TestConvertBodyProcessor(t *testing.T) {
	p := NewConvertBody(reflect.TypeOf(0), nil, &string2int{})
	e := exchange.NewExchange(context.TODO())
	e.Message().SetBody("333")

	p.Process(e)

	assert.Equal(t, 333, e.MessageBody())
}
