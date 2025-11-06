package convertheader

import (
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
	"reflect"
)

type converter interface {
	Convert(any, reflect.Type, map[string]any) (any, error)
}

type convertHeaderProcessor struct {
	routeName  string
	name       string
	headerName string
	params     map[string]any
	targetType reflect.Type
	conv       converter
}

func NewProcessor(routeName, name, headerName string, targetType reflect.Type, params map[string]any, conv converter) *convertHeaderProcessor {
	return &convertHeaderProcessor{
		routeName:  routeName,
		name:       name,
		headerName: headerName,
		params:     params,
		targetType: targetType,
		conv:       conv,
	}
}

func (p *convertHeaderProcessor) Name() string {
	return p.name
}

func (p *convertHeaderProcessor) RouteName() string {
	return p.routeName
}

func (p *convertHeaderProcessor) Process(e *exchange.Exchange) {
	if headerValue, headerExists := e.Message().Header(p.headerName); headerExists {
		convertedHeader, err := p.conv.Convert(headerValue, p.targetType, p.params)
		if err != nil {
			e.SetError(err)
			return
		}

		e.Message().SetHeader(p.headerName, convertedHeader)
	}
}
