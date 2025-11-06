package convertbody

import (
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
	"reflect"
)

type converter interface {
	Convert(any, reflect.Type, map[string]any) (any, error)
}

type convertBodyProcessor struct {
	routeName  string
	name       string
	params     map[string]any
	targetType reflect.Type
	conv       converter
}

func NewProcessor(routeName, name string, targetType reflect.Type, params map[string]any, conv converter) *convertBodyProcessor {
	return &convertBodyProcessor{
		routeName:  routeName,
		name:       name,
		params:     params,
		targetType: targetType,
		conv:       conv,
	}
}

func (p *convertBodyProcessor) Name() string {
	return p.name
}

func (p *convertBodyProcessor) RouteName() string {
	return p.routeName
}

func (p *convertBodyProcessor) Process(e *exchange.Exchange) {
	convertedBody, err := p.conv.Convert(e.Message().Body, p.targetType, p.params)
	if err != nil {
		e.SetError(err)
		return
	}

	e.Message().Body = convertedBody
}
