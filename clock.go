package logger

import "time"

// DefaultClock uses internal Go time-handling functions.
type DefaultClock struct {
}

func NewDefaultClock() *DefaultClock {
	return &DefaultClock{}
}

// Provides current time same as time.Now().
func (c *DefaultClock) Now() time.Time {
	return time.Now()
}
