package processor

import (
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
	"reflect"
)

type ConvertHeader struct {
	headerName string
	targetType reflect.Type
	params     map[string]any
	conv       Converter
}

func NewConvertHeader(name string, targetType reflect.Type, params map[string]any, conv Converter) *ConvertHeader {
	return &ConvertHeader{
		headerName: name,
		params:     params,
		targetType: targetType,
		conv:       conv,
	}
}

func (p *ConvertHeader) Process(e *exchange.Exchange) {
	if value, exists := e.Message().Header(p.headerName); exists {
		convertedValue, err := p.conv.Convert(value, p.targetType, p.params)
		if err != nil {
			e.SetError(err)
			return
		}

		e.SetMessageHeader(p.headerName, convertedValue)
	}
}
