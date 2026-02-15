package converter

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIntToString(t *testing.T) {
	s, err := IntToString().Convert(123, nil)
	assert.NoError(t, err)
	assert.Equal(t, "123", s)
}
