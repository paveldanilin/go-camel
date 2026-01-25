package processor

import (
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
	"reflect"
)

type ConvertProperty struct {
	propertyName string
	targetType   reflect.Type
	params       map[string]any
	conv         Converter
}

func NewConvertProperty(name string, targetType reflect.Type, params map[string]any, conv Converter) *ConvertProperty {
	return &ConvertProperty{
		propertyName: name,
		targetType:   targetType,
		params:       params,
		conv:         conv,
	}
}

func (p *ConvertProperty) Process(e *exchange.Exchange) {
	if value, exists := e.Property(p.propertyName); exists {
		convertedValue, err := p.conv.Convert(value, p.targetType, p.params)
		if err != nil {
			e.SetError(err)
			return
		}

		e.SetProperty(p.propertyName, convertedValue)
	}
}
