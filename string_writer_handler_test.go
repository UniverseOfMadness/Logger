package logger

import (
	"bytes"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestStringWriterHandler_Handle(t *testing.T) {
	t.Parallel()

	writer := &bytes.Buffer{}
	handler := NewStringWriterHandler(writer)

	err := handler.Handle(Log{
		Level:     LevelDebug,
		Message:   "test msg",
		Data:      make(Data),
		CreatedAt: time.Now(),
	})

	if assert.NoError(t, err) {
		assert.Equal(t, "test msg\n", writer.String())
	}
}

func TestStringWriterHandler_Handle_WithFormatter(t *testing.T) {
	t.Parallel()

	log := Log{
		Level:     LevelDebug,
		Message:   "test msg",
		Data:      make(Data),
		CreatedAt: time.Now(),
	}

	mFormatter := &mockFormatter{}
	mFormatter.On("Format", log).Return(FormattedLog{
		Log:              log,
		FormattedMessage: "formatted test msg",
	})

	writer := &bytes.Buffer{}

	handler := NewStringWriterHandler(writer)
	handler.UseFormatter(mFormatter)

	err := handler.Handle(log)

	if assert.NoError(t, err) {
		assert.Equal(t, "formatted test msg\n", writer.String())
		mFormatter.AssertExpectations(t)
	}
}

func TestStringWriterHandler_Handle_Failure(t *testing.T) {
	t.Parallel()

	mWriter := &mockStringWriter{}
	mWriter.On("WriteString", "test msg\n").Return(0, errors.New("test"))

	handler := NewStringWriterHandler(mWriter)

	err := handler.Handle(Log{
		Level:     LevelDebug,
		Message:   "test msg",
		Data:      make(Data),
		CreatedAt: time.Now(),
	})

	if assert.Error(t, err) {
		assert.EqualError(t, err, "StringWriterHandler - error occurred while handling log: test")
	}
}

func TestStringWriterHandler_HandleBatch(t *testing.T) {
	t.Parallel()

	writer := &bytes.Buffer{}
	handler := NewStringWriterHandler(writer)

	err := handler.HandleBatch([]Log{
		{
			Level:     LevelDebug,
			Message:   "test msg",
			Data:      make(Data),
			CreatedAt: time.Now(),
		},
		{
			Level:     LevelError,
			Message:   "test err",
			Data:      make(Data),
			CreatedAt: time.Now(),
		},
	})

	if assert.NoError(t, err) {
		assert.Equal(t, "test msg\ntest err\n", writer.String())
	}
}

func TestStringWriterHandler_HandleBatch_WithFormatter(t *testing.T) {
	t.Parallel()

	log := Log{
		Level:     LevelDebug,
		Message:   "test msg",
		Data:      make(Data),
		CreatedAt: time.Now(),
	}

	mFormatter := &mockFormatter{}
	mFormatter.On("Format", mock.IsType(Log{})).Return(FormattedLog{
		Log:              log,
		FormattedMessage: "formatted test msg",
	})

	writer := &bytes.Buffer{}
	handler := NewStringWriterHandler(writer)
	handler.UseFormatter(mFormatter)

	err := handler.HandleBatch([]Log{
		{
			Level:     LevelDebug,
			Message:   "test msg",
			Data:      make(Data),
			CreatedAt: time.Now(),
		},
		{
			Level:     LevelError,
			Message:   "test err",
			Data:      make(Data),
			CreatedAt: time.Now(),
		},
	})

	if assert.NoError(t, err) {
		assert.Equal(t, "formatted test msg\nformatted test msg\n", writer.String())
	}
}

func TestStringWriterHandler_HandleBatch_Failure(t *testing.T) {
	t.Parallel()

	mWriter := &mockStringWriter{}
	mWriter.On("WriteString", "test msg\ntest err\n").Return(0, errors.New("test"))

	handler := NewStringWriterHandler(mWriter)

	err := handler.HandleBatch([]Log{
		{
			Level:     LevelDebug,
			Message:   "test msg",
			Data:      make(Data),
			CreatedAt: time.Now(),
		},
		{
			Level:     LevelError,
			Message:   "test err",
			Data:      make(Data),
			CreatedAt: time.Now(),
		},
	})

	if assert.Error(t, err) {
		assert.EqualError(t, err, "StringWriterHandler - error occurred while handling logs: test")
	}
}
