package convertproperty

import (
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
	"reflect"
)

type converter interface {
	Convert(any, reflect.Type, map[string]any) (any, error)
}

type convertPropertyProcessor struct {
	routeName    string
	name         string
	propertyName string
	params       map[string]any
	targetType   reflect.Type
	conv         converter
}

func NewProcessor(routeName, name, propertyName string, targetType reflect.Type, params map[string]any, conv converter) *convertPropertyProcessor {
	return &convertPropertyProcessor{
		routeName:    routeName,
		name:         name,
		propertyName: propertyName,
		params:       params,
		targetType:   targetType,
		conv:         conv,
	}
}

func (p *convertPropertyProcessor) Name() string {
	return p.name
}

func (p *convertPropertyProcessor) RouteName() string {
	return p.routeName
}

func (p *convertPropertyProcessor) Process(e *exchange.Exchange) {
	if propertyValue, propertyExists := e.Property(p.propertyName); propertyExists {
		convertedProperty, err := p.conv.Convert(propertyValue, p.targetType, p.params)
		if err != nil {
			e.SetError(err)
			return
		}

		e.SetProperty(p.propertyName, convertedProperty)
	}
}
