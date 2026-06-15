package logger

import (
	"fmt"

	"github.com/TheRootDaemon/tlgc/termcolor"
)

// log formats and writes a message if the level is enabled.
func (l *Logger) log(level Level, format string, args []any) {
	if !l.Enabled(level) {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	prefix := termcolor.Sprint(levelColors[level], levelLabels[level])
	msg := fmt.Sprintf(format, args...)
	_, _ = fmt.Fprintf(l.writer, "%s: %s\n", prefix, msg)
}

// Trace logs a message at LevelTrace.
func (l *Logger) Trace(format string, args ...any) {
	l.log(LevelTrace, format, args)
}

// Debug logs a message at LevelDebug.
func (l *Logger) Debug(format string, args ...any) {
	l.log(LevelDebug, format, args)
}

// Info logs a message at LevelInfo.
func (l *Logger) Info(format string, args ...any) {
	l.log(LevelInfo, format, args)
}

// Warn logs a message at LevelWarn.
func (l *Logger) Warn(format string, args ...any) {
	l.log(LevelWarn, format, args)
}

// Error logs a message at LevelError.
func (l *Logger) Error(format string, args ...any) {
	l.log(LevelError, format, args)
}

// Log logs a message at the given level.
func (l *Logger) Log(level Level, format string, args ...any) {
	l.log(level, format, args)
}

// InfoStart writes a status message without a trailing newline.
//
// In debug and trace modes it falls back to Info.
// Use InfoEnd to complete the line.
func (l *Logger) InfoStart(format string, args ...any) {
	if !l.Enabled(LevelInfo) {
		return
	}

	if l.Enabled(LevelDebug) {
		l.Info(format, args...)
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	prefix := termcolor.Sprint(levelColors[LevelInfo], levelLabels[LevelInfo])
	msg := fmt.Sprintf(format, args...)
	_, _ = fmt.Fprintf(l.writer, "%s %s", prefix, msg)
}

// InfoEnd writes the trailing part of a status message
// started with InfoStart.
//
// It is a no-op in debug and trace modes
// (where InfoStart already completed the line).
func (l *Logger) InfoEnd(format string, args ...any) {
	if !l.Enabled(LevelInfo) || l.Enabled(LevelDebug) {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	msg := fmt.Sprintf(format, args...)
	_, _ = fmt.Fprintf(l.writer, "%s\n", msg)
}
