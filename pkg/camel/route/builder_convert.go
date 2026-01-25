package route

import "maps"

type ConvertBuilder struct {
	builder *Builder
}

type BodyConvertBuilder struct {
	builder    *Builder
	targetType any
	typeName   string
	params     map[string]any
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
	if (bc.targetType == nil && len(bc.typeName) == 0) {
		// TODO: error
	}

	if (bc.targetType != nil) {
		return bc.builder.ConvertBodyTo(bc.targetType, bc.params)
	}

	return bc.builder.ConvertBodyToNamed(bc.typeName, bc.params)
}

func (cb *ConvertBuilder) Body() *BodyConvertBuilder {
	return &BodyConvertBuilder{builder: cb.builder}
}

func (cb *ConvertBuilder) Header(name string) {

}

func (cb *ConvertBuilder) Property(name string) {

}

func (cb *ConvertBuilder) Variable(name string) {

}