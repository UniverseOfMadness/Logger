package logger

import "time"

type (
	// CriticalHandleFunc is called when Logger.Critical or MainLogger.Criticalf.
	CriticalHandleFunc func(message string, data Data)
	// FailureHandleFunc is called when handler returns an error.
	FailureHandleFunc func(log Log, err error)
	// Clock provides time for logs.
	// Log.CreatedAt value in Log will be created
	// using Now function.
	Clock interface {
		// Provides current time.
		Now() time.Time
	}
	// Handlers must be designed to process incoming
	// logs and store them / execute actions related to them.
	// Handle function needs to accept log prepared by main
	// *MainLogger instance.
	Handler interface {
		// Handle processes given log message according to
		// its implementation.
		Handle(log Log) error
	}
	// Formatters creates final presentation of message.
	// Formatter must provide log with prepared FormattedLog.FormattedMessage.
	// Handler will decide which version of message will be used.
	Formatter interface {
		// Format decorates provided Log with
		// FormattedLog which includes FormattedLog.FormattedMessage
		// field containing specially formatted message for handlers.
		Format(log Log) FormattedLog
	}
	// DebugLogger contains only functions for debug log messages.
	DebugLogger interface {
		// Debug creates Log with LevelDebug and provided values as
		// Data in Log. Each "key:value" pair will be saved separately,
		// for example Debug("test", "key", "val", "x", "y") will create
		// Data type with {"key":"val","x":"y"}. Values count for Data must be even.
		Debug(message string, values ...string)
		// Debugf creates Log with LevelDebug and formats message with
		// fmt.Sprintf function with all values before passing it to handler.
		Debugf(message string, values ...interface{})
	}
	// InfoLogger contains only functions for info log messages.
	InfoLogger interface {
		// Info creates Log with LevelInfo and provided values as
		// Data in Log. Each "key:value" pair will be saved separately,
		// for example Info("test", "key", "val", "x", "y") will create
		// Data type with {"key":"val","x":"y"}. Values count for Data must be even.
		Info(message string, values ...string)
		// Infof creates Log with LevelInfo and formats message with
		// fmt.Sprintf function with all values before passing it to handler.
		Infof(message string, values ...interface{})
	}
	// WarningLogger contains only functions for warning log messages.
	WarningLogger interface {
		// Warning creates Log with LevelWarning and provided values as
		// Data in Log. Each "key:value" pair will be saved separately,
		// for example Warning("test", "key", "val", "x", "y") will create
		// Data type with {"key":"val","x":"y"}. Values count for Data must be even.
		Warning(message string, values ...string)
		// Warningf creates Log with LevelWarning and formats message with
		// fmt.Sprintf function with all values before passing it to handler.
		Warningf(message string, values ...interface{})
	}
	// ErrorLogger contains only functions for error log messages.
	ErrorLogger interface {
		// Error creates Log with LevelError and provided values as
		// Data in Log. Each "key:value" pair will be saved separately,
		// for example Error("test", "key", "val", "x", "y") will create
		// Data type with {"key":"val","x":"y"}. Values count for Data must be even.
		Error(message string, values ...string)
		// Errorf creates Log with LevelError and formats message with
		// fmt.Sprintf function with all values before passing it to handler.
		Errorf(message string, values ...interface{})
	}
	// CriticalLogger contains only functions for error log messages.
	CriticalLogger interface {
		// Critical creates Log with LevelCritical and provided values as
		// Data in Log. Each "key:value" pair will be saved separately,
		// for example Critical("test", "key", "val", "x", "y") will create
		// Data type with {"key":"val","x":"y"}. Values count for Data must be even.
		Critical(message string, values ...string)
		// Criticalf creates Log with LevelCritical and formats message with
		// fmt.Sprintf function with all values before passing it to handler.
		Criticalf(message string, values ...interface{})
	}
	// Logger is aggregator for all logger sub-interfaces:
	//  * DebugLogger
	//  * InfoLogger
	//  * WarningLogger
	//  * ErrorLogger
	//  * CriticalLogger
	Logger interface {
		DebugLogger
		InfoLogger
		WarningLogger
		ErrorLogger
		CriticalLogger
	}
)
