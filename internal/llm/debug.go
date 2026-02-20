package llm

import (
	"fmt"
	"os"
	"time"
)

// DebugLogger provides structured debug logging with timestamps.
// Output goes to stderr and can be controlled via DEBUG=1 environment variable.
type DebugLogger struct {
	name    string
	enabled bool
}

// NewDebugLogger creates a new debug logger. It's enabled if DEBUG=1 is set.
func NewDebugLogger(name string) *DebugLogger {
	enabled := os.Getenv("DEBUG") == "1"
	return &DebugLogger{
		name:    name,
		enabled: enabled,
	}
}

// Log logs a message with timestamp and context name.
// Format: [HH:MM:SS.mmm] [context] message
func (dl *DebugLogger) Log(format string, args ...interface{}) {
	if !dl.enabled {
		return
	}
	msg := fmt.Sprintf(format, args...)
	timestamp := time.Now().Format("15:04:05.000")
	fmt.Fprintf(os.Stderr, "[%s] [%s] %s\n", timestamp, dl.name, msg)
}

// LogError logs an error with context.
func (dl *DebugLogger) LogError(context string, err error) {
	if !dl.enabled {
		return
	}
	timestamp := time.Now().Format("15:04:05.000")
	fmt.Fprintf(os.Stderr, "[%s] [%s] ERROR: %s: %v\n", timestamp, dl.name, context, err)
}
