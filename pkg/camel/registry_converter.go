package camel

import "reflect"

type convertChainStep struct {
	toType reflect.Type
	conv   any
}

// convertChainPath - cache for cache of converters.
type convertChainPath struct {
	steps []convertChainStep
}

type converterRegistry struct {
	converters map[reflect.Type]map[reflect.Type]any // map[fromType][toType]Converter[From, To]
	cache      map[string]convertChainPath           // cacheKey: "from.String():to.String()" -> chainPath
	named      map[string]reflect.Type
}

func NewConverterRegistry() ConverterRegistry {
	return &converterRegistry{
		converters: make(map[reflect.Type]map[reflect.Type]any),
		cache:      make(map[string]convertChainPath),
		named:      make(map[string]reflect.Type),
	}
}
