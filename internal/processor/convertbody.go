package processor

import (
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
	"reflect"
)

type bodyConverter interface {
	Convert(any, reflect.Type, map[string]any) (any, error)
}

type ConvertBody struct {
	targetType reflect.Type
	params     map[string]any
	conv       bodyConverter
}

func NewConvertBody(targetType reflect.Type, params map[string]any, conv bodyConverter) *ConvertBody {
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

	e.Message().SetBody(convertedValue)
}
