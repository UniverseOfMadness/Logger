package logger

import "time"

// CriticalHandleFunc is called when Logger.Critical or Logger.Criticalf.
type CriticalHandleFunc func(message string, data Data)

// FailureHandleFunc is called when handler returns an error.
type FailureHandleFunc func(log Log, err error)

// Clock provides time for logs.
// Log.CreatedAt value in Log will be created
// using Now function.
type Clock interface {
	// Provides current time.
	Now() time.Time
}

// Handlers must be designed to process incoming
// logs and store them / execute actions related to them.
// Handle function needs to accept log prepared by main
// *Logger instance.
type Handler interface {
	// Handle processes given log message according to
	// its implementation.
	Handle(log Log) error
}

// Formatters creates final presentation of message.
// Formatter must provide log with prepared FormattedLog.FormattedMessage.
// Handler will decide which version of message will be used.
type Formatter interface {
	// Format decorates provided Log with
	// FormattedLog which includes FormattedLog.FormattedMessage
	// field containing specially formatted message for handlers.
	Format(log Log) FormattedLog
}
