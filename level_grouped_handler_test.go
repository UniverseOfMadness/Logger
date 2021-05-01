package logger

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestLevelGroupedHandler_Handle(t *testing.T) {
	t.Parallel()

	logInfo := Log{
		Level:     LevelInfo,
		Message:   "test info",
		Data:      make(Data),
		CreatedAt: time.Now(),
	}

	logWarning := Log{
		Level:     LevelWarning,
		Message:   "test warning",
		Data:      make(Data),
		CreatedAt: time.Now(),
	}

	logError := Log{
		Level:     LevelError,
		Message:   "test error",
		Data:      make(Data),
		CreatedAt: time.Now(),
	}

	mFallbackHandler := &mockHandler{}
	mFallbackHandler.On("Handle", logWarning).Return(nil)

	mHandler1 := &mockHandler{}
	mHandler1.On("Handle", logInfo).Return(nil)

	mHandler2 := &mockHandler{}
	mHandler2.On("Handle", logError).Return(nil)

	handler := NewLevelGroupedHandler(mFallbackHandler, LevelGroup{
		Levels:  []Level{LevelInfo},
		Handler: mHandler1,
	}, LevelGroup{
		Levels:  []Level{LevelError},
		Handler: mHandler2,
	})

	infoErr := handler.Handle(logInfo)
	assert.NoError(t, infoErr)

	warningErr := handler.Handle(logWarning)
	assert.NoError(t, warningErr)

	errorErr := handler.Handle(logError)

	if assert.NoError(t, errorErr) {
		mFallbackHandler.AssertExpectations(t)
		mHandler1.AssertExpectations(t)
		mHandler2.AssertExpectations(t)
	}
}

func TestLevelGroupedHandler_HandleBatch(t *testing.T) {
	t.Parallel()

	logs := []Log{
		{
			Level:     LevelInfo,
			Message:   "test info",
			Data:      make(Data),
			CreatedAt: time.Now(),
		},
		{
			Level:     LevelWarning,
			Message:   "test warning",
			Data:      make(Data),
			CreatedAt: time.Now(),
		},
		{
			Level:     LevelError,
			Message:   "test error",
			Data:      make(Data),
			CreatedAt: time.Now(),
		},
	}

	mFallbackHandler := &mockHandler{}
	mFallbackHandler.On("HandleBatch", []Log{logs[1]}).Return(nil)

	mHandler1 := &mockHandler{}
	mHandler1.On("HandleBatch", []Log{logs[0]}).Return(nil)

	mHandler2 := &mockHandler{}
	mHandler2.On("HandleBatch", []Log{logs[2]}).Return(nil)

	handler := NewLevelGroupedHandler(mFallbackHandler, LevelGroup{
		Levels:  []Level{LevelInfo},
		Handler: mHandler1,
	}, LevelGroup{
		Levels:  []Level{LevelError},
		Handler: mHandler2,
	})

	hErr := handler.HandleBatch(logs)

	if assert.NoError(t, hErr) {
		mFallbackHandler.AssertExpectations(t)
		mHandler1.AssertExpectations(t)
		mHandler2.AssertExpectations(t)
	}
}

func TestLevelGroupedHandler_Handle_FallbackHandlerFailure(t *testing.T) {
	t.Parallel()

	logInfo := Log{
		Level:     LevelInfo,
		Message:   "test info",
		Data:      make(Data),
		CreatedAt: time.Now(),
	}

	logWarning := Log{
		Level:     LevelWarning,
		Message:   "test warning",
		Data:      make(Data),
		CreatedAt: time.Now(),
	}

	logError := Log{
		Level:     LevelError,
		Message:   "test error",
		Data:      make(Data),
		CreatedAt: time.Now(),
	}

	mFallbackHandler := &mockHandler{}
	mFallbackHandler.On("Handle", logWarning).Return(errors.New("test"))

	mHandler1 := &mockHandler{}
	mHandler1.On("Handle", logInfo).Return(nil)

	mHandler2 := &mockHandler{}
	mHandler2.On("Handle", logError).Return(nil)

	handler := NewLevelGroupedHandler(mFallbackHandler, LevelGroup{
		Levels:  []Level{LevelInfo},
		Handler: mHandler1,
	}, LevelGroup{
		Levels:  []Level{LevelError},
		Handler: mHandler2,
	})

	infoErr := handler.Handle(logInfo)
	assert.NoError(t, infoErr)

	errorErr := handler.Handle(logError)
	assert.NoError(t, errorErr)

	warningErr := handler.Handle(logWarning)

	if assert.Error(t, warningErr) {
		assert.EqualError(t, warningErr, "LevelGroupedHandler - fallback handler returned an error: test")

		mFallbackHandler.AssertExpectations(t)
		mHandler1.AssertExpectations(t)
		mHandler2.AssertExpectations(t)
	}
}

func TestLevelGroupedHandler_Handle_GroupedHandlerFailure(t *testing.T) {
	t.Parallel()

	logInfo := Log{
		Level:     LevelInfo,
		Message:   "test info",
		Data:      make(Data),
		CreatedAt: time.Now(),
	}

	logWarning := Log{
		Level:     LevelWarning,
		Message:   "test warning",
		Data:      make(Data),
		CreatedAt: time.Now(),
	}

	logError := Log{
		Level:     LevelError,
		Message:   "test error",
		Data:      make(Data),
		CreatedAt: time.Now(),
	}

	mFallbackHandler := &mockHandler{}
	mFallbackHandler.On("Handle", logWarning).Return(nil)

	mHandler1 := &mockHandler{}
	mHandler1.On("Handle", logInfo).Return(errors.New("test"))

	mHandler2 := &mockHandler{}
	mHandler2.On("Handle", logError).Return(nil)

	handler := NewLevelGroupedHandler(mFallbackHandler, LevelGroup{
		Levels:  []Level{LevelInfo},
		Handler: mHandler1,
	}, LevelGroup{
		Levels:  []Level{LevelError},
		Handler: mHandler2,
	})

	errorErr := handler.Handle(logError)
	assert.NoError(t, errorErr)

	warningErr := handler.Handle(logWarning)
	assert.NoError(t, warningErr)

	infoErr := handler.Handle(logInfo)

	if assert.Error(t, infoErr) {
		assert.EqualError(t, infoErr, "LevelGroupedHandler - one of handlers returned an error: test")

		mFallbackHandler.AssertExpectations(t)
		mHandler1.AssertExpectations(t)
		mHandler2.AssertExpectations(t)
	}
}

func TestLevelGroupedHandler_HandleBatch_FallbackHandlerFailure(t *testing.T) {
	t.Parallel()

	logs := []Log{
		{
			Level:     LevelInfo,
			Message:   "test info",
			Data:      make(Data),
			CreatedAt: time.Now(),
		},
		{
			Level:     LevelWarning,
			Message:   "test warning",
			Data:      make(Data),
			CreatedAt: time.Now(),
		},
		{
			Level:     LevelError,
			Message:   "test error",
			Data:      make(Data),
			CreatedAt: time.Now(),
		},
	}

	mFallbackHandler := &mockHandler{}
	mFallbackHandler.On("HandleBatch", []Log{logs[1]}).Return(errors.New("test"))

	mHandler1 := &mockHandler{}
	mHandler1.On("HandleBatch", []Log{logs[0]}).Return(nil)

	mHandler2 := &mockHandler{}

	handler := NewLevelGroupedHandler(mFallbackHandler, LevelGroup{
		Levels:  []Level{LevelInfo},
		Handler: mHandler1,
	}, LevelGroup{
		Levels:  []Level{LevelError},
		Handler: mHandler2,
	})

	hErr := handler.HandleBatch(logs)

	if assert.Error(t, hErr) {
		assert.EqualError(t, hErr, "LevelGroupedHandler - fallback handler returned an error: test")

		mFallbackHandler.AssertExpectations(t)
		mHandler1.AssertExpectations(t)
		mHandler2.AssertNotCalled(t, "HandleBatch", []Log{logs[2]})
	}
}

func TestLevelGroupedHandler_HandleBatch_GroupedHandlerFailure(t *testing.T) {
	t.Parallel()

	logs := []Log{
		{
			Level:     LevelInfo,
			Message:   "test info",
			Data:      make(Data),
			CreatedAt: time.Now(),
		},
		{
			Level:     LevelWarning,
			Message:   "test warning",
			Data:      make(Data),
			CreatedAt: time.Now(),
		},
		{
			Level:     LevelError,
			Message:   "test error",
			Data:      make(Data),
			CreatedAt: time.Now(),
		},
	}

	mFallbackHandler := &mockHandler{}
	mFallbackHandler.On("HandleBatch", []Log{logs[1]}).Return(nil)

	mHandler1 := &mockHandler{}
	mHandler1.On("HandleBatch", []Log{logs[0]}).Return(nil)

	mHandler2 := &mockHandler{}
	mHandler2.On("HandleBatch", []Log{logs[2]}).Return(errors.New("test"))

	handler := NewLevelGroupedHandler(mFallbackHandler, LevelGroup{
		Levels:  []Level{LevelInfo},
		Handler: mHandler1,
	}, LevelGroup{
		Levels:  []Level{LevelError},
		Handler: mHandler2,
	})

	hErr := handler.HandleBatch(logs)

	if assert.Error(t, hErr) {
		assert.EqualError(t, hErr, "LevelGroupedHandler - one of handlers returned an error: test")

		mFallbackHandler.AssertExpectations(t)
		mHandler1.AssertExpectations(t)
		mHandler2.AssertExpectations(t)
	}
}
