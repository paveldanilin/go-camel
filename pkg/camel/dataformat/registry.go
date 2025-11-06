package dataformat

import (
	"errors"
	"github.com/paveldanilin/go-camel/pkg/camel/api"
	"sync"
)

type registry struct {
	mu            sync.Mutex
	dataFormatMap map[string]api.DataFormat
}

func NewRegistry() *registry {
	return &registry{
		dataFormatMap: map[string]api.DataFormat{},
	}
}

func (r *registry) RegisterDataFormat(name string, dataFormat api.DataFormat) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.dataFormatMap[name]; exists {
		return errors.New("data format already registered")
	}

	r.dataFormatMap[name] = dataFormat
	return nil
}

func (r *registry) DataFormat(name string) api.DataFormat {
	r.mu.Lock()
	defer r.mu.Unlock()

	if df, exists := r.dataFormatMap[name]; exists {
		return df
	}

	return nil
}
