package converter

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFloat64ToString(t *testing.T) {
	s, err := Float64ToString().Convert(1.12345, nil)
	assert.NoError(t, err)
	assert.Equal(t, "1.12345", s)
}

func TestFloat64ToString_Format(t *testing.T) {
	s, err := Float64ToString().Convert(1.12345, map[string]any{
		FloatFormat: byte('e'),
	})
	assert.NoError(t, err)
	assert.Equal(t, "1.12345e+00", s)
}
