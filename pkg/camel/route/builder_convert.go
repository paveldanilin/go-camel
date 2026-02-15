package route

import (
	"github.com/paveldanilin/go-camel/pkg/camel/route/step"
	"maps"
)

type ConvertBuilder struct {
	builder *Builder
	opts    []step.Option
}

type BodyConvertBuilder struct {
	routeBuilder   *Builder
	convertBuilder *ConvertBuilder
	targetType     any
	typeName       string
	params         map[string]any
}

func (bc *BodyConvertBuilder) To(targetType any) *BodyConvertBuilder {
	bc.targetType = targetType
	return bc
}

func (bc *BodyConvertBuilder) ToType(typeName string) *BodyConvertBuilder {
	bc.typeName = typeName
	return bc
}

func (bc *BodyConvertBuilder) WithParams(params map[string]any) *BodyConvertBuilder {
	maps.Copy(bc.params, params)
	return bc
}

func (bc *BodyConvertBuilder) EndConvert() *Builder {
	if bc.targetType == nil && len(bc.typeName) == 0 {
		// TODO: error
	}

	if bc.targetType != nil {
		return bc.routeBuilder.ConvertBodyTo(bc.targetType, bc.params, bc.convertBuilder.opts...)
	}

	return bc.routeBuilder.ConvertBodyToNamed(bc.typeName, bc.params, bc.convertBuilder.opts...)
}

func (cb *ConvertBuilder) Body() *BodyConvertBuilder {
	return &BodyConvertBuilder{
		routeBuilder:   cb.builder,
		convertBuilder: cb,
	}
}

func (cb *ConvertBuilder) Header(name string) {

}

func (cb *ConvertBuilder) Property(name string) {

}

func (cb *ConvertBuilder) Variable(name string) {

}
