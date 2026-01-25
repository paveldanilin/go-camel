package processor

import (
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
	"reflect"
)

type ConvertVariable struct {
	variableName string
	targetType   reflect.Type
	params       map[string]any
	conv         Converter
}

func NewConvertVariable(name string, targetType reflect.Type, params map[string]any, conv Converter) *ConvertVariable {
	return &ConvertVariable{
		variableName: name,
		targetType:   targetType,
		params:       params,
		conv:         conv,
	}
}

func (p *ConvertVariable) Process(e *exchange.Exchange) {
	if value, exists := e.Variable(p.variableName); exists {
		convertedValue, err := p.conv.Convert(value, p.targetType, p.params)
		if err != nil {
			e.SetError(err)
			return
		}

		e.SetVariable(p.variableName, convertedValue)
	}
}
