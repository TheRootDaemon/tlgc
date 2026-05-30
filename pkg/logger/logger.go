package logger

import (
	"io"
	"os"
	"sync"
)

// Level specifies the security level of a log.
type Level int

const (
	LevelTrace Level = iota - 2
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
)

// levelColors maps log levels to terminal color styles.
var levelColors = map[Level]string{
	LevelTrace: "blue bold",
	LevelDebug: "magenta bold",
	LevelInfo:  "cyan bold",
	LevelWarn:  "yellow bold",
	LevelError: "red bold",
}

// levelLabels maps log levels to their display labels.
var levelLabels = map[Level]string{
	LevelTrace: "trace:",
	LevelDebug: "debug:",
	LevelInfo:  "info:",
	LevelWarn:  "warning:",
	LevelError: "error:",
}

// Logger writes leveled log messages to a writer.
type Logger struct {
	mu     sync.Mutex
	level  Level
	writer io.Writer
}

// New creates a Logger that writes to stderr.
//
// Quiet suppresses all output below LevelError.
// Each level of verbose enables the next lower level
func New(quiet bool, verbose uint8) *Logger {
	return NewWithWriter(quiet, verbose, os.Stderr)
}

// NewWithWriter creates a Logger with the given writer.
func NewWithWriter(quiet bool, verbose uint8, writer io.Writer) *Logger {
	level := LevelInfo

	switch {
	case quiet:
		level = LevelError
	case verbose >= 2:
		level = LevelTrace
	case verbose >= 1:
		level = LevelDebug
	}

	return &Logger{
		level:  level,
		writer: writer,
	}
}

// Enabled reports whether messages at the given level  would be logged.
func (l *Logger) Enabled(level Level) bool {
	return l.level <= level
}
