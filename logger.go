package logger

import "fmt"

// MainLogger struct wraps handler "Handle" function
// into few simple functions for easier log management.
type MainLogger struct {
	handler         Handler
	clock           Clock
	criticalHandler CriticalHandleFunc
	failureHandler  FailureHandleFunc
	config          *config
}

// New creates new *MainLogger instance.
func New(handler Handler) *MainLogger {
	return &MainLogger{handler: handler, clock: NewDefaultClock(), config: newConfig(LevelDebug)}
}

// SetLevel changes minimum required Level for Log to be handled (LevelDebug by default - all logs).
func (l *MainLogger) SetLevel(level Level) {
	l.config.setLevel(level)
}

// WithClock allows to set custom implementation for
// Clock interface and manage log time with it.
func (l *MainLogger) WithClock(clock Clock) *MainLogger {
	l.clock = clock

	return l
}

// WithCriticalHandler sets custom handler for Critical and Criticalf logs.
// By default no handler is set.
func (l *MainLogger) WithCriticalHandler(handleFunc CriticalHandleFunc) *MainLogger {
	l.criticalHandler = handleFunc

	return l
}

func (l *MainLogger) WithFailureHandler(handleFunc FailureHandleFunc) *MainLogger {
	l.failureHandler = handleFunc

	return l
}

func (l *MainLogger) Debug(message string, values ...string) {
	l.handleStandardLog(LevelDebug, message, values)
}

func (l *MainLogger) Debugf(message string, values ...interface{}) {
	l.handleFormattedLog(LevelDebug, message, values)
}

func (l *MainLogger) Info(message string, values ...string) {
	l.handleStandardLog(LevelInfo, message, values)
}

func (l *MainLogger) Infof(message string, values ...interface{}) {
	l.handleFormattedLog(LevelInfo, message, values)
}

func (l *MainLogger) Warning(message string, values ...string) {
	l.handleStandardLog(LevelWarning, message, values)
}

func (l *MainLogger) Warningf(message string, values ...interface{}) {
	l.handleFormattedLog(LevelWarning, message, values)
}

func (l *MainLogger) Error(message string, values ...string) {
	l.handleStandardLog(LevelError, message, values)
}

func (l *MainLogger) Errorf(message string, values ...interface{}) {
	l.handleFormattedLog(LevelError, message, values)
}

func (l *MainLogger) Critical(message string, values ...string) {
	l.handleWithCritical(l.handleStandardLog(LevelCritical, message, values))
}

func (l *MainLogger) Criticalf(message string, values ...interface{}) {
	l.handleWithCritical(l.handleFormattedLog(LevelCritical, message, values))
}

func (l *MainLogger) handleStandardLog(level Level, message string, values []string) (Log, bool) {
	if !level.EqualOrGreaterThan(l.config.getLevel()) {
		return Log{}, false
	}

	log := l.createLog(level, message, values)
	l.handleError(log, l.handler.Handle(log))

	return log, true
}

func (l *MainLogger) handleFormattedLog(level Level, message string, values []interface{}) (Log, bool) {
	if !level.EqualOrGreaterThan(l.config.getLevel()) {
		return Log{}, false
	}

	log := l.createLog(level, fmt.Sprintf(message, values...), []string{})
	l.handleError(log, l.handler.Handle(log))

	return log, true
}

func (l *MainLogger) handleWithCritical(log Log, isHandling bool) {
	if isHandling && l.criticalHandler != nil {
		l.criticalHandler(log.Message, log.Data)
	}
}

func (l *MainLogger) handleError(log Log, err error) {
	if l.failureHandler != nil {
		l.failureHandler(log, err)
	}
}

func (l *MainLogger) createLog(level Level, message string, values []string) Log {
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
