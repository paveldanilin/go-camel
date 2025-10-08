package camel

import (
	"errors"
	"sync"
)

type DataFormat interface {
	Unmarshal(data []byte, targetType any) (any, error)
	Marshal(data any) error
}

// ---------------------------------------------------------------------------------------------------------------------
// DataFormatRegistry
// ---------------------------------------------------------------------------------------------------------------------

type DataFormatRegistry interface {
	RegisterDataFormat(name string, format DataFormat) error
	DataFormat(name string) DataFormat
}

type dataFormatRegistry struct {
	mu            sync.Mutex
	dataFormatMap map[string]DataFormat
}

func newDataFormatRegistry() *dataFormatRegistry {
	return &dataFormatRegistry{
		dataFormatMap: map[string]DataFormat{},
	}
}

func (r *dataFormatRegistry) RegisterDataFormat(name string, dataFormat DataFormat) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.dataFormatMap[name]; exists {
		return errors.New("data format already registered")
	}

	r.dataFormatMap[name] = dataFormat
	return nil
}

func (r *dataFormatRegistry) DataFormat(name string) DataFormat {
	r.mu.Lock()
	defer r.mu.Unlock()

	if df, exists := r.dataFormatMap[name]; exists {
		return df
	}

	return nil
}
