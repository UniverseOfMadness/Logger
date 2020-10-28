package logger

import (
	"github.com/stretchr/testify/mock"
	"os"
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

type mockFilesystem struct {
	mock.Mock
}

func (m *mockFilesystem) OpenFile(name string, flag int, perm os.FileMode) (File, error) {
	args := m.Called(name, flag, perm)
	return args.Get(0).(File), args.Error(1)
}

type mockFile struct {
	mock.Mock
}

func (m *mockFile) WriteString(s string) (n int, err error) {
	args := m.Called(s)
	return args.Int(0), args.Error(1)
}

func (m *mockFile) Write(p []byte) (n int, err error) {
	args := m.Called(p)
	return args.Int(0), args.Error(1)
}

func (m *mockFile) Close() error {
	return m.Called().Error(0)
}
