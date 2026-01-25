package processor

import (
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
	"reflect"
)

type ConvertBody struct {
	targetType reflect.Type
	params     map[string]any
	conv       Converter
}

func NewConvertBody(targetType reflect.Type, params map[string]any, conv Converter) *ConvertBody {
	return &ConvertBody{
		targetType: targetType,
		params:     params,
		conv:       conv,
	}
}

func (p *ConvertBody) Process(e *exchange.Exchange) {
	convertedValue, err := p.conv.Convert(e.Message().Body(), p.targetType, p.params)
	if err != nil {
		e.SetError(err)
		return
	}

	e.SetMessageBody(convertedValue)
}
