package api

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

// LogLevel represents the severity of a log message.
type LogLevel string

const (
	LogInfo    LogLevel = "info"
	LogWarn    LogLevel = "warn"
	LogError   LogLevel = "error"
	LogDebug   LogLevel = "debug"
)

// StructuredLogger outputs JSON-formatted log lines for production observability.
// Falls back to standard log.Printf for human-readable dev output.
type StructuredLogger struct {
	structured bool
	component  string
}

// NewStructuredLogger creates a logger. If LOG_FORMAT=json, outputs structured JSON.
func NewStructuredLogger(component string) *StructuredLogger {
	format := os.Getenv("LOG_FORMAT")
	return &StructuredLogger{
		structured: format == "json",
		component:  component,
	}
}

type logEntry struct {
	Timestamp string   `json:"timestamp"`
	Level     LogLevel `json:"level"`
	Component string   `json:"component"`
	Message   string   `json:"message"`
	Data      any      `json:"data,omitempty"`
}

func (l *StructuredLogger) log(level LogLevel, msg string, data any) {
	if l.structured {
		entry := logEntry{
			Timestamp: time.Now().UTC().Format(time.RFC3339Nano),
			Level:     level,
			Component: l.component,
			Message:   msg,
			Data:      data,
		}
		out, _ := json.Marshal(entry)
		fmt.Fprintln(os.Stdout, string(out))
	} else {
		prefix := ""
		switch level {
		case LogInfo:
			prefix = "ℹ️"
		case LogWarn:
			prefix = "⚠️"
		case LogError:
			prefix = "❌"
		case LogDebug:
			prefix = "🔍"
		}
		if data != nil {
			log.Printf("%s [%s] %s | %+v", prefix, l.component, msg, data)
		} else {
			log.Printf("%s [%s] %s", prefix, l.component, msg)
		}
	}
}

func (l *StructuredLogger) Info(msg string, data ...any) {
	var d any
	if len(data) > 0 {
		d = data[0]
	}
	l.log(LogInfo, msg, d)
}

func (l *StructuredLogger) Warn(msg string, data ...any) {
	var d any
	if len(data) > 0 {
		d = data[0]
	}
	l.log(LogWarn, msg, d)
}

func (l *StructuredLogger) Error(msg string, data ...any) {
	var d any
	if len(data) > 0 {
		d = data[0]
	}
	l.log(LogError, msg, d)
}

func (l *StructuredLogger) Debug(msg string, data ...any) {
	var d any
	if len(data) > 0 {
		d = data[0]
	}
	l.log(LogDebug, msg, d)
}
