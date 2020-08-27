package logger

import "fmt"

// Logger struct wraps handler "Handle" function
// into few simple functions for easier log management.
type Logger struct {
	handler         Handler
	clock           Clock
	criticalHandler CriticalHandleFunc
	failureHandler  FailureHandleFunc
	config          *config
}

// New creates new *Logger instance.
func New(handler Handler) *Logger {
	return &Logger{handler: handler, clock: NewDefaultClock(), config: newConfig(LevelDebug)}
}

// SetLevel changes minimum required Level for Log to be handled (LevelDebug by default - all logs).
func (l *Logger) SetLevel(level Level) {
	l.config.setLevel(level)
}

// WithClock allows to set custom implementation for
// Clock interface and manage log time with it.
func (l *Logger) WithClock(clock Clock) *Logger {
	l.clock = clock

	return l
}

// WithCriticalHandler sets custom handler for Critical and Criticalf logs.
// By default no handler is set.
func (l *Logger) WithCriticalHandler(handleFunc CriticalHandleFunc) *Logger {
	l.criticalHandler = handleFunc

	return l
}

func (l *Logger) WithFailureHandler(handleFunc FailureHandleFunc) *Logger {
	l.failureHandler = handleFunc

	return l
}

// Debug creates Log with LevelDebug and provided values as
// Data in Log. Each "key:value" pair will be saved separately,
// for example Debug("test", "key", "val", "x", "y") will create
// Data type with {"key":"val","x":"y"}. Values count for Data must be even.
func (l *Logger) Debug(message string, values ...string) {
	l.handleStandardLog(LevelDebug, message, values)
}

// Debugf creates Log with LevelDebug and formats message with
// fmt.Sprintf function with all values before passing it to handler.
func (l *Logger) Debugf(message string, values ...interface{}) {
	l.handleFormattedLog(LevelDebug, message, values)
}

// Info creates Log with LevelInfo and provided values as
// Data in Log. Each "key:value" pair will be saved separately,
// for example Info("test", "key", "val", "x", "y") will create
// Data type with {"key":"val","x":"y"}. Values count for Data must be even.
func (l *Logger) Info(message string, values ...string) {
	l.handleStandardLog(LevelInfo, message, values)
}

// Infof creates Log with LevelInfo and formats message with
// fmt.Sprintf function with all values before passing it to handler.
func (l *Logger) Infof(message string, values ...interface{}) {
	l.handleFormattedLog(LevelInfo, message, values)
}

// Warning creates Log with LevelWarning and provided values as
// Data in Log. Each "key:value" pair will be saved separately,
// for example Warning("test", "key", "val", "x", "y") will create
// Data type with {"key":"val","x":"y"}. Values count for Data must be even.
func (l *Logger) Warning(message string, values ...string) {
	l.handleStandardLog(LevelWarning, message, values)
}

// Warningf creates Log with LevelWarning and formats message with
// fmt.Sprintf function with all values before passing it to handler.
func (l *Logger) Warningf(message string, values ...interface{}) {
	l.handleFormattedLog(LevelWarning, message, values)
}

// Error creates Log with LevelError and provided values as
// Data in Log. Each "key:value" pair will be saved separately,
// for example Error("test", "key", "val", "x", "y") will create
// Data type with {"key":"val","x":"y"}. Values count for Data must be even.
func (l *Logger) Error(message string, values ...string) {
	l.handleStandardLog(LevelError, message, values)
}

// Errorf creates Log with LevelError and formats message with
// fmt.Sprintf function with all values before passing it to handler.
func (l *Logger) Errorf(message string, values ...interface{}) {
	l.handleFormattedLog(LevelError, message, values)
}

// Critical creates Log with LevelCritical and provided values as
// Data in Log. Each "key:value" pair will be saved separately,
// for example Critical("test", "key", "val", "x", "y") will create
// Data type with {"key":"val","x":"y"}. Values count for Data must be even.
func (l *Logger) Critical(message string, values ...string) {
	l.handleWithCritical(l.handleStandardLog(LevelCritical, message, values))
}

// Criticalf creates Log with LevelCritical and formats message with
// fmt.Sprintf function with all values before passing it to handler.
func (l *Logger) Criticalf(message string, values ...interface{}) {
	l.handleWithCritical(l.handleFormattedLog(LevelCritical, message, values))
}

func (l *Logger) handleStandardLog(level Level, message string, values []string) (Log, bool) {
	if !level.EqualOrGreaterThan(l.config.getLevel()) {
		return Log{}, false
	}

	log := l.createLog(level, message, values)
	l.handleError(log, l.handler.Handle(log))

	return log, true
}

func (l *Logger) handleFormattedLog(level Level, message string, values []interface{}) (Log, bool) {
	if !level.EqualOrGreaterThan(l.config.getLevel()) {
		return Log{}, false
	}

	log := l.createLog(level, fmt.Sprintf(message, values...), []string{})
	l.handleError(log, l.handler.Handle(log))

	return log, true
}

func (l *Logger) handleWithCritical(log Log, isHandling bool) {
	if isHandling && l.criticalHandler != nil {
		l.criticalHandler(log.Message, log.Data)
	}
}

func (l *Logger) handleError(log Log, err error) {
	if l.failureHandler != nil {
		l.failureHandler(log, err)
	}
}

func (l *Logger) createLog(level Level, message string, values []string) Log {
	d := make(Data)

	if len(values) > 0 {
		d = createDataFromSlice(values)
	}

	return Log{
		Level:     level,
		Message:   message,
		Data:      d,
		CreatedAt: l.clock.Now(),
	}
}
