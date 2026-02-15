package converter

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStringToDateTime(t *testing.T) {
	c := StringToDateTime()

	time, err := c.Convert("2010-05-05 05:05:05", nil)

	assert.NoError(t, err)
	assert.Equal(t, 2010, time.Year())
	assert.Equal(t, 5, int(time.Month()))
	assert.Equal(t, 5, time.Day())
}
