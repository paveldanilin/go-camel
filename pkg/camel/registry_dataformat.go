package camel

import (
	"errors"
	"github.com/paveldanilin/go-camel/pkg/camel/api"
	"sync"
)

type DataFormatRegistry interface {
	RegisterDataFormat(name string, format api.DataFormat) error
	DataFormat(name string) api.DataFormat
}

type dataFormatRegistry struct {
	mu            sync.Mutex
	dataFormatMap map[string]api.DataFormat
}

func newDataFormatRegistry() *dataFormatRegistry {
	return &dataFormatRegistry{
		dataFormatMap: map[string]api.DataFormat{},
	}
}

func (r *dataFormatRegistry) RegisterDataFormat(name string, dataFormat api.DataFormat) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.dataFormatMap[name]; exists {
		return errors.New("data format already registered")
	}

	r.dataFormatMap[name] = dataFormat
	return nil
}

func (r *dataFormatRegistry) DataFormat(name string) api.DataFormat {
	r.mu.Lock()
	defer r.mu.Unlock()

	if df, exists := r.dataFormatMap[name]; exists {
		return df
	}

	return nil
}
