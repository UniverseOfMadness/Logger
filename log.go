package logger

import "time"

type Log struct {
	Level     Level
	Message   string
	Data      Data
	CreatedAt time.Time
}

type FormattedLog struct {
	Log
	FormattedMessage string
}
