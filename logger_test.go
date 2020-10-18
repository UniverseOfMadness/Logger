package logger

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestLogger_Log(t *testing.T) {
	t.Parallel()

	message := "test message"
	messageWithFormatting := "test %s message"
	additionalEvenValues := []string{"key", "val"}
	additionalOddValues := []string{"key"}
	expectedMessage := "test message"
	expectedFormattedMessage := "test formatted message"
	expectedData := Data{"key": "val"}

	t.Run("debug", func(t *testing.T) {
		testLoggerLogWithoutAdditionalValues(t, func(logger *MainLogger) {
			logger.Debug(message)
		}, expectedMessage, LevelDebug)

		testLoggerLogWithEvenNumberOfAdditionalValues(t, func(logger *MainLogger) {
			logger.Debug(message, additionalEvenValues...)
		}, expectedMessage, LevelDebug, expectedData)

		testLoggerLogWithOddNumberOfAdditionalValues(t, func(logger *MainLogger) {
			logger.Debug(message, additionalOddValues...)
		})

		testLoggerLogfWithFormatting(t, func(logger *MainLogger) {
			logger.Debugf(messageWithFormatting, "formatted")
		}, expectedFormattedMessage, LevelDebug)
	})

	t.Run("info", func(t *testing.T) {
		testLoggerLogWithoutAdditionalValues(t, func(logger *MainLogger) {
			logger.Info(message)
		}, expectedMessage, LevelInfo)

		testLoggerLogWithEvenNumberOfAdditionalValues(t, func(logger *MainLogger) {
			logger.Info(message, additionalEvenValues...)
		}, expectedMessage, LevelInfo, expectedData)

		testLoggerLogWithOddNumberOfAdditionalValues(t, func(logger *MainLogger) {
			logger.Info(message, additionalOddValues...)
		})

		testLoggerLogfWithFormatting(t, func(logger *MainLogger) {
			logger.Infof(messageWithFormatting, "formatted")
		}, expectedFormattedMessage, LevelInfo)
	})

	t.Run("warning", func(t *testing.T) {
		testLoggerLogWithoutAdditionalValues(t, func(logger *MainLogger) {
			logger.Warning(message)
		}, expectedMessage, LevelWarning)

		testLoggerLogWithEvenNumberOfAdditionalValues(t, func(logger *MainLogger) {
			logger.Warning(message, additionalEvenValues...)
		}, expectedMessage, LevelWarning, expectedData)

		testLoggerLogWithOddNumberOfAdditionalValues(t, func(logger *MainLogger) {
			logger.Warning(message, additionalOddValues...)
		})

		testLoggerLogfWithFormatting(t, func(logger *MainLogger) {
			logger.Warningf(messageWithFormatting, "formatted")
		}, expectedFormattedMessage, LevelWarning)
	})

	t.Run("error", func(t *testing.T) {
		testLoggerLogWithoutAdditionalValues(t, func(logger *MainLogger) {
			logger.Error(message)
		}, expectedMessage, LevelError)

		testLoggerLogWithEvenNumberOfAdditionalValues(t, func(logger *MainLogger) {
			logger.Error(message, additionalEvenValues...)
		}, expectedMessage, LevelError, expectedData)

		testLoggerLogWithOddNumberOfAdditionalValues(t, func(logger *MainLogger) {
			logger.Error(message, additionalOddValues...)
		})

		testLoggerLogfWithFormatting(t, func(logger *MainLogger) {
			logger.Errorf(messageWithFormatting, "formatted")
		}, expectedFormattedMessage, LevelError)
	})

	t.Run("critical", func(t *testing.T) {
		testLoggerLogWithoutAdditionalValues(t, func(logger *MainLogger) {
			logger.Critical(message)
		}, expectedMessage, LevelCritical)

		testLoggerLogWithEvenNumberOfAdditionalValues(t, func(logger *MainLogger) {
			logger.Critical(message, additionalEvenValues...)
		}, expectedMessage, LevelCritical, expectedData)

		testLoggerLogWithOddNumberOfAdditionalValues(t, func(logger *MainLogger) {
			logger.Critical(message, additionalOddValues...)
		})

		testLoggerLogfWithFormatting(t, func(logger *MainLogger) {
			logger.Criticalf(messageWithFormatting, "formatted")
		}, expectedFormattedMessage, LevelCritical)
	})
}

func TestLogger_Critical_WithCustomHandler(t *testing.T) {
	t.Parallel()

	res := ""
	mHandler := &mockHandler{}
	mHandler.On("Handle", mock.IsType(Log{})).Return(nil)

	logger := New(mHandler)
	logger.WithCriticalHandler(func(message string, _ Data) {
		res = message
	})

	logger.Critical("test message")

	assert.NotEmpty(t, res)
	assert.Equal(t, "test message", res)

	mHandler.AssertExpectations(t)
}

func TestLogger_Criticalf_WithCustomHandler(t *testing.T) {
	t.Parallel()

	res := ""
	mHandler := &mockHandler{}
	mHandler.On("Handle", mock.IsType(Log{})).Return(nil)

	logger := New(mHandler)
	logger.WithCriticalHandler(func(message string, _ Data) {
		res = message
	})

	logger.Criticalf("test %s message", "formatted")

	assert.NotEmpty(t, res)
	assert.Equal(t, "test formatted message", res)

	mHandler.AssertExpectations(t)
}

func TestLogger_Log_WithHandlerFailure(t *testing.T) {
	t.Parallel()

	executed := false
	tm := time.Now()

	expectedLog := Log{
		Level:     LevelWarning,
		Message:   "test message",
		Data:      make(Data),
		CreatedAt: tm,
	}

	mClock := &mockClock{}
	mClock.On("Now").Return(tm)

	mHandler := &mockHandler{}
	mHandler.On("Handle", expectedLog).Return(errors.New("test"))

	logger := New(mHandler)
	logger.WithClock(mClock)
	logger.WithFailureHandler(func(log Log, err error) {
		if assert.Error(t, err) {
			assert.Equal(t, expectedLog, log)
			assert.EqualError(t, err, "test")

			executed = true
		}
	})

	logger.Warning("test message")

	mHandler.AssertExpectations(t)
	mClock.AssertExpectations(t)

	assert.True(t, executed)
}

func TestLogger_Log_WithSetLevel(t *testing.T) {
	t.Parallel()

	tm := time.Now()

	logInfo := Log{
		Level:     LevelInfo,
		Message:   "test info",
		Data:      make(Data),
		CreatedAt: tm,
	}

	logError := Log{
		Level:     LevelError,
		Message:   "test error",
		Data:      make(Data),
		CreatedAt: tm,
	}

	mClock := &mockClock{}
	mClock.On("Now").Return(tm)

	mHandler := &mockHandler{}
	mHandler.AssertNotCalled(t, "Handle", logInfo)
	mHandler.On("Handle", logError).Return(nil)

	logger := New(mHandler)
	logger.SetLevel(LevelWarning)
	logger.WithClock(mClock)

	logger.Info("test info")
	logger.Infof("test info")

	logger.Error("test error")
	logger.Errorf("test error")

	mHandler.AssertExpectations(t)
	mHandler.AssertNumberOfCalls(t, "Handle", 2)

	mClock.AssertExpectations(t)
}

func testLoggerLogWithoutAdditionalValues(
	t *testing.T,
	loggerCallback func(logger *MainLogger),
	expectedMessage string,
	expectedLevel Level,
) {
	tm := time.Now()

	mClock := &mockClock{}
	mClock.On("Now").Return(tm)

	mHandler := &mockHandler{}
	mHandler.On("Handle", Log{
		Level:     expectedLevel,
		Message:   expectedMessage,
		Data:      Data{},
		CreatedAt: tm,
	}).Return(nil)

	logger := New(mHandler)
	logger.WithClock(mClock)

	loggerCallback(logger)

	mClock.AssertExpectations(t)
	mHandler.AssertExpectations(t)
}

func testLoggerLogWithEvenNumberOfAdditionalValues(
	t *testing.T,
	loggerCallback func(logger *MainLogger),
	expectedMessage string,
	expectedLevel Level,
	expectedData Data,
) {
	tm := time.Now()

	mClock := &mockClock{}
	mClock.On("Now").Return(tm)

	mHandler := &mockHandler{}
	mHandler.On("Handle", Log{
		Level:     expectedLevel,
		Message:   expectedMessage,
		Data:      expectedData,
		CreatedAt: tm,
	}).Return(nil)

	logger := New(mHandler)
	logger.WithClock(mClock)

	loggerCallback(logger)

	mClock.AssertExpectations(t)
	mHandler.AssertExpectations(t)
}

func testLoggerLogWithOddNumberOfAdditionalValues(t *testing.T, loggerCallback func(logger *MainLogger)) {
	mClock := &mockClock{}
	mHandler := &mockHandler{}

	logger := New(mHandler)
	logger.WithClock(mClock)

	assert.PanicsWithValue(t, "number of items in slice provided for Data must be even", func() {
		loggerCallback(logger)
	})

	mClock.AssertNotCalled(t, "Now")
	mHandler.AssertNotCalled(t, "Handle")
}

func testLoggerLogfWithFormatting(
	t *testing.T,
	loggerCallback func(logger *MainLogger),
	expectedMessage string,
	expectedLevel Level,
) {
	tm := time.Now()

	mClock := &mockClock{}
	mClock.On("Now").Return(tm)

	mHandler := &mockHandler{}
	mHandler.On("Handle", Log{
		Level:     expectedLevel,
		Message:   expectedMessage,
		Data:      Data{},
		CreatedAt: tm,
	}).Return(nil)

	logger := New(mHandler)
	logger.WithClock(mClock)

	loggerCallback(logger)

	mClock.AssertExpectations(t)
	mHandler.AssertExpectations(t)
}
