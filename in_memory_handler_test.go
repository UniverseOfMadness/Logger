package logger

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestInMemoryHandler_Handle(t *testing.T) {
	t.Parallel()

	log := Log{
		Level:     LevelDebug,
		Message:   "test",
		Data:      make(Data),
		CreatedAt: time.Now(),
	}

	handler := NewInMemoryHandler(0)
	err := handler.Handle(log)

	if assert.NoError(t, err) {
		assert.False(t, handler.IsEmpty())
		assert.Equal(t, log, handler.Pop())
		assert.True(t, handler.IsEmpty())
	}
}

func TestInMemoryHandler_Pop(t *testing.T) {
	t.Parallel()

	log := Log{
		Level:     LevelDebug,
		Message:   "test",
		Data:      make(Data),
		CreatedAt: time.Now(),
	}

	handler := NewInMemoryHandler(0)
	err := handler.Handle(log)

	if assert.NoError(t, err) {
		assert.Equal(t, log, handler.Pop())

		assert.PanicsWithValue(t, "no logs in handler", func() {
			handler.Pop()
		})
	}
}

func TestInMemoryHandler_Clear(t *testing.T) {
	t.Parallel()

	log := Log{
		Level:     LevelDebug,
		Message:   "test",
		Data:      make(Data),
		CreatedAt: time.Now(),
	}

	handler := NewInMemoryHandler(0)

	_ = handler.Handle(log)
	_ = handler.Handle(log)
	_ = handler.Handle(log)
	_ = handler.Handle(log)

	assert.False(t, handler.IsEmpty())

	handler.Clear()
	assert.True(t, handler.IsEmpty())
}

func TestInMemoryHandler_BufferOverflow(t *testing.T) {
	t.Parallel()

	log := Log{
		Level:     LevelDebug,
		Message:   "test",
		Data:      make(Data),
		CreatedAt: time.Now(),
	}

	handler := NewInMemoryHandler(2)

	err1 := handler.Handle(log)
	assert.NoError(t, err1)

	err2 := handler.Handle(log)
	assert.NoError(t, err2)

	err3 := handler.Handle(log)
	assert.Error(t, err3)
	assert.EqualError(t, err3, "InMemoryHandler - number of logs exceeded buffer limit (2 records)")
}
