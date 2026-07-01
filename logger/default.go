package logger

import "os"

// defaultLogger is the package-level Logger used by the
// convenience functions.
var defaultLogger = New(false, 0)

// exit is a package-level variable so tests can replace it.
var exit = os.Exit

// SetDefault sets the package-level default logger.
func SetDefault(l *Logger) {
	defaultLogger = l
}

// Trace logs at LevelTrace via the default logger.
func Trace(format string, args ...any) {
	defaultLogger.Trace(format, args...)
}

// Debug logs at LevelDebug via the default logger.
func Debug(format string, args ...any) {
	defaultLogger.Debug(format, args...)
}

// Info logs at LevelInfo via the default logger.
func Info(format string, args ...any) {
	defaultLogger.Info(format, args...)
}

// Warn logs at LevelWarn via the default logger.
func Warn(format string, args ...any) {
	defaultLogger.Warn(format, args...)
}

// Error logs at LevelError via the default logger.
func Error(format string, args ...any) {
	defaultLogger.Error(format, args...)
}

// Exit logs the message at LevelError via the default logger
// and terminates the program with exit status 1.
func Exit(format string, args ...any) {
	defaultLogger.Error(format, args...)
	exit(1)
}

// Log logs at the given level via the default logger.
func Log(level Level, format string, args ...any) {
	defaultLogger.Log(level, format, args...)
}

// InfoStart writes a status line start via the default logger.
func InfoStart(format string, args ...any) {
	defaultLogger.InfoStart(format, args...)
}

// InfoEnd completes a status line via the default logger.
func InfoEnd(format string, args ...any) {
	defaultLogger.InfoEnd(format, args...)
}

// Enabled reports whether the given level is enabled
// on the default logger.
func Enabled(level Level) bool {
	return defaultLogger.Enabled(level)
}
