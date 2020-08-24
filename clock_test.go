package logger

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestDefaultClock_Now(t *testing.T) {
	t.Parallel()

	clock := NewDefaultClock()
	tm := clock.Now()

	assert.IsType(t, time.Time{}, tm)
	assert.Equal(t, time.Now().Format(time.RFC1123Z), tm.Format(time.RFC1123Z))
}
