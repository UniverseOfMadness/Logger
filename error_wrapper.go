package logger

import "fmt"

type ErrorWrappedLogger struct {
	*Logger
}

func NewErrorWrappedLogger(logger *Logger) *ErrorWrappedLogger {
	return &ErrorWrappedLogger{Logger: logger}
}

func (w *ErrorWrappedLogger) OnError(err error, values ...string) {
	w.on(w.Logger.Error, err, values...)
}

func (w *ErrorWrappedLogger) OnErrorWrapped(err error, message string, values ...interface{}) {
	w.onWrapped(w.Logger.Error, err, message, values...)
}

func (w *ErrorWrappedLogger) OnCritical(err error, values ...string) {
	w.on(w.Logger.Critical, err, values...)
}

func (w *ErrorWrappedLogger) OnCriticalWrapped(err error, message string, values ...interface{}) {
	w.onWrapped(w.Logger.Critical, err, message, values...)
}

func (w *ErrorWrappedLogger) on(call func(message string, values ...string), err error, values ...string) {
	if err != nil {
		call(err.Error(), values...)
	}
}

func (w *ErrorWrappedLogger) onWrapped(call func(message string, values ...string), err error, message string, values ...interface{}) {
	if err != nil {
		values = append(values, err)

		call(fmt.Errorf(message, values...).Error())
	}
}
