package processor

import (
	"context"
	"errors"
	"github.com/paveldanilin/go-camel/pkg/camel/env"
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
	"github.com/paveldanilin/go-camel/pkg/camel/template"
	camelUri "github.com/paveldanilin/go-camel/pkg/camel/uri"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

type endpointProducerFactory struct {
}

func (endpointProducerFactory) CreateEndpointProducer(uri string) (Processor, error) {
	if strings.HasPrefix(uri, "http://srv.api.com/services") {
		return P(func(e *exchange.Exchange) {
			parsedUri, _ := camelUri.Parse(uri, nil)
			if parsedUri.MustParam("key") == "qwerty1234567890" {
				e.SetMessageBody(`{"code": "OK", "result": "SERVICE_ACTIVATED", "service": "` + parsedUri.MustParam("serviceCode") + `"}`)
			} else {
				e.SetMessageBody(`{"code": "ERROR", "result": "BAD_KEY"}`)
			}
		}), nil
	}
	return nil, errors.New("unknown uri")
}

func TestToD(t *testing.T) {
	uriTpl, err := template.Parse("http://srv.api.com/services?key=${env.key}&serviceCode=${header.service_code}")
	assert.NoError(t, err)

	tod := NewToD(uriTpl, uriTpl.Vars(), &endpointProducerFactory{}, env.NewMap(map[string]string{
		"key": "qwerty1234567890",
	}))

	// First
	ex := exchange.NewExchange(context.TODO())
	ex.SetMessageHeader("service_code", "TEL1")

	tod.Process(ex)

	assert.NoError(t, ex.Error())
	assert.Equal(t, `{"code": "OK", "result": "SERVICE_ACTIVATED", "service": "TEL1"}`, ex.MessageBody())

	// Second
	ex2 := exchange.NewExchange(context.TODO())
	ex2.SetMessageHeader("service_code", "BLK1")

	tod.Process(ex2)

	assert.NoError(t, ex2.Error())
	assert.Equal(t, `{"code": "OK", "result": "SERVICE_ACTIVATED", "service": "BLK1"}`, ex2.MessageBody())
}
