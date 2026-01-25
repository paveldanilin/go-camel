package dataformat

import (
	"errors"
	"sync"
)

type Registry struct {
	mu            sync.Mutex
	dataFormatMap map[string]DataFormat
}

func NewRegistry() *Registry {
	return &Registry{
		dataFormatMap: map[string]DataFormat{},
	}
}

func (r *Registry) RegisterDataFormat(name string, dataFormat DataFormat) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.dataFormatMap[name]; exists {
		return errors.New("data format already registered")
	}

	r.dataFormatMap[name] = dataFormat
	return nil
}

func (r *Registry) DataFormat(name string) DataFormat {
	r.mu.Lock()
	defer r.mu.Unlock()

	if df, exists := r.dataFormatMap[name]; exists {
		return df
	}

	return nil
}
