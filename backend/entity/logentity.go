package entity

import "time"

type LogLevel string

const (
	LogLevelInfo  LogLevel = "INFO"
	LogLevelError LogLevel = "ERROR"
	LogLevelPanic LogLevel = "PANIC"
)

type LogEntry struct {
	Level     LogLevel
	Timestamp time.Time
	Source    string
	Message   string
}
