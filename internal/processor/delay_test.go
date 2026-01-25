package processor

import (
	"context"
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestDelayProcessor(t *testing.T) {
	p := NewDelay(500)
	e := exchange.NewExchange(context.TODO())
	start := time.Now()

	p.Process(e)

	elapsedMs := time.Since(start).Milliseconds()

	assert.LessOrEqual(t, int64(500), elapsedMs)
}
