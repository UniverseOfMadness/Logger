package logger

import (
	"errors"
	"testing"
	"time"
)

func TestErrorWrappedLogger_OnError(t *testing.T) {
	t.Parallel()

	log := Log{
		Level:     LevelError,
		Message:   "this is an {text}",
		Data:      Data{"text": "error"},
		CreatedAt: time.Now(),
	}

	mHandler := &mockHandler{}
	mHandler.On("Handle", log).Return(nil)

	mClock := &mockClock{}
	mClock.On("Now").Return(log.CreatedAt)

	logger := New(mHandler)
	logger.WithClock(mClock)
	wrapper := NewErrorWrappedLogger(logger)

	wrapper.OnError(errors.New("this is an {text}"), "text", "error")

	mHandler.AssertExpectations(t)
	mClock.AssertExpectations(t)
}

func TestErrorWrappedLogger_OnErrorWrapped(t *testing.T) {
	t.Parallel()

	log := Log{
		Level:     LevelError,
		Message:   "something important went wrong: the error",
		Data:      make(Data),
		CreatedAt: time.Now(),
	}

	mHandler := &mockHandler{}
	mHandler.On("Handle", log).Return(nil)

	mClock := &mockClock{}
	mClock.On("Now").Return(log.CreatedAt)

	logger := New(mHandler)
	logger.WithClock(mClock)
	wrapper := NewErrorWrappedLogger(logger)

	wrapper.OnErrorWrapped(errors.New("the error"), "something %s went wrong: %w", "important")

	mHandler.AssertExpectations(t)
	mClock.AssertExpectations(t)
}

func TestErrorWrappedLogger_OnCritical(t *testing.T) {
	t.Parallel()

	log := Log{
		Level:     LevelCritical,
		Message:   "this is an {text}",
		Data:      Data{"text": "error"},
		CreatedAt: time.Now(),
	}

	mHandler := &mockHandler{}
	mHandler.On("Handle", log).Return(nil)

	mClock := &mockClock{}
	mClock.On("Now").Return(log.CreatedAt)

	logger := New(mHandler)
	logger.WithClock(mClock)
	wrapper := NewErrorWrappedLogger(logger)

	wrapper.OnCritical(errors.New("this is an {text}"), "text", "error")

	mHandler.AssertExpectations(t)
	mClock.AssertExpectations(t)
}

func TestErrorWrappedLogger_OnCriticalWrapped(t *testing.T) {
	t.Parallel()

	log := Log{
		Level:     LevelCritical,
		Message:   "something important went wrong: the error",
		Data:      make(Data),
		CreatedAt: time.Now(),
	}

	mHandler := &mockHandler{}
	mHandler.On("Handle", log).Return(nil)

	mClock := &mockClock{}
	mClock.On("Now").Return(log.CreatedAt)

	logger := New(mHandler)
	logger.WithClock(mClock)
	wrapper := NewErrorWrappedLogger(logger)

	wrapper.OnCriticalWrapped(errors.New("the error"), "something %s went wrong: %w", "important")

	mHandler.AssertExpectations(t)
	mClock.AssertExpectations(t)
}
