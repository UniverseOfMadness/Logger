package logger

import (
	"github.com/stretchr/testify/mock"
	"time"
)

type mockClock struct {
	mock.Mock
}

func (c *mockClock) Now() time.Time {
	return c.Called().Get(0).(time.Time)
}

type mockHandler struct {
	mock.Mock
}

func (h *mockHandler) HandleBatch(logs []Log) error {
	return h.Called(logs).Error(0)
}

func (h *mockHandler) Handle(log Log) error {
	return h.Called(log).Error(0)
}

type mockFormatter struct {
	mock.Mock
}

func (f *mockFormatter) Format(log Log) FormattedLog {
	return f.Called(log).Get(0).(FormattedLog)
}

type mockStringWriter struct {
	mock.Mock
}

func (w *mockStringWriter) WriteString(s string) (n int, err error) {
	args := w.Called(s)

	return args.Int(0), args.Error(1)
}
