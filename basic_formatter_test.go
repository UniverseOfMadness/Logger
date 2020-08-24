package logger

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestBasicFormatter_Format(t *testing.T) {
	t.Parallel()

	t.Run("without additional data", func(t *testing.T) {
		tm := time.Now()

		formatter := NewBasicFormatter("testing", time.RFC3339)
		formatted := formatter.Format(Log{
			Level:     LevelDebug,
			Message:   "test message",
			Data:      Data{},
			CreatedAt: tm,
		})

		assert.IsType(t, FormattedLog{}, formatted)
		assert.Equal(t, fmt.Sprintf("testing | %s | DEBUG | test message", tm.Format(time.RFC3339)), formatted.FormattedMessage)
	})

	t.Run("with additional data", func(t *testing.T) {
		tm := time.Now()

		formatter := NewBasicFormatter("testing", time.RFC3339)
		formatted := formatter.Format(Log{
			Level:     LevelDebug,
			Message:   "test message {answer}",
			Data:      Data{"key": "val", "answer": "42"},
			CreatedAt: tm,
		})

		assert.IsType(t, FormattedLog{}, formatted)

		assert.Equal(
			t,
			fmt.Sprintf("testing | %s | DEBUG | test message 42 | answer:42 key:val", tm.Format(time.RFC3339)),
			formatted.FormattedMessage,
		)
	})
}
