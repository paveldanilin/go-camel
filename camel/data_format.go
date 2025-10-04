package camel

import (
	"errors"
	"sync"
)

type DataFormat interface {
	Read(data []byte, targetType any) (any, error)
	Write(data any) error
}

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

func (dr *dataFormatRegistry) RegisterDataFormat(name string, dataFormat DataFormat) error {
	dr.mu.Lock()
	defer dr.mu.Unlock()

	if _, exists := dr.dataFormatMap[name]; exists {
		return errors.New("data format already registered")
	}

	dr.dataFormatMap[name] = dataFormat
	return nil
}

func (dr *dataFormatRegistry) DataFormat(name string) DataFormat {
	dr.mu.Lock()
	defer dr.mu.Unlock()

	if df, exists := dr.dataFormatMap[name]; exists {
		return df
	}

	return nil
}
