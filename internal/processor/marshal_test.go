package processor

import (
	"context"
	"github.com/paveldanilin/go-camel/pkg/camel/dataformat"
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
	"github.com/stretchr/testify/assert"
	"testing"
)

type jsonModelMarshal struct {
	Page   int
	Fruits []string
}

func TestJsonProcessorMarshal(t *testing.T) {
	p := NewMarshal(&dataformat.JSON{})

	e := exchange.NewExchange(context.TODO())
	e.Message().SetBody(jsonModelMarshal{Page: 2, Fruits: []string{"orange", "abricot"}})

	p.Process(e)

	assert.False(t, e.IsError())
	assert.Equal(t, `{"Page":2,"Fruits":["orange","abricot"]}`, e.MessageBody())
}
