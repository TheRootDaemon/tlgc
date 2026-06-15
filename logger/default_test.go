package logger

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultLogger(t *testing.T) {
	assert.NotNil(t, defaultLogger)
	assert.Equal(t, LevelInfo, defaultLogger.level)
}

func TestSetDefault(t *testing.T) {
	var buf bytes.Buffer
	l := &Logger{
		level:  LevelInfo,
		writer: &buf,
	}

	old := defaultLogger
	t.Cleanup(func() { defaultLogger = old })

	SetDefault(l)
	Info("hello %s", "world")

	output := buf.String()
	assert.Contains(t, output, "INFO")
	assert.Contains(t, output, "hello world")
}

func TestPackageLevelFunctions(t *testing.T) {
	var buf bytes.Buffer
	l := &Logger{
		level:  LevelTrace, // enable every level
		writer: &buf,
	}

	old := defaultLogger
	t.Cleanup(func() { defaultLogger = old })
	SetDefault(l)

	tests := []struct {
		name    string
		logFunc func(string, ...any)
		label   string
	}{
		{name: "Trace", logFunc: Trace, label: "TRACE"},
		{name: "Debug", logFunc: Debug, label: "DEBUG"},
		{name: "Info", logFunc: Info, label: "INFO"},
		{name: "Warn", logFunc: Warn, label: "WARNING"},
		{name: "Error", logFunc: Error, label: "ERROR"},
		{
			name: "Log",
			logFunc: func(format string, args ...any) {
				Log(LevelInfo, format, args...)
			},
			label: "INFO",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()
			tt.logFunc("hello %s", "world")
			output := buf.String()
			assert.Contains(t, output, tt.label)
			assert.Contains(t, output, "hello world")
		})
	}
}

func TestPackageLevelEnabled(t *testing.T) {
	tests := []struct {
		name     string
		level    Level // default logger level
		arg      Level // level to query
		expected bool
	}{
		{name: "info when error", level: LevelError, arg: LevelInfo, expected: false},
		{name: "error when error", level: LevelError, arg: LevelError, expected: true},
		{name: "trace when info", level: LevelInfo, arg: LevelTrace, expected: false},
		{name: "info when info", level: LevelInfo, arg: LevelInfo, expected: true},
		{name: "warn when info", level: LevelInfo, arg: LevelWarn, expected: true},
		{name: "trace when trace", level: LevelTrace, arg: LevelDebug, expected: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repl := &Logger{level: tt.level}

			old := defaultLogger
			t.Cleanup(func() { defaultLogger = old })
			SetDefault(repl)

			assert.Equal(t, tt.expected, Enabled(tt.arg))
		})
	}
}

func TestExit(t *testing.T) {
	var buf bytes.Buffer
	l := &Logger{
		level:  LevelError,
		writer: &buf,
	}

	oldLogger := defaultLogger
	t.Cleanup(func() { defaultLogger = oldLogger })
	SetDefault(l)

	var exitCode int
	oldExit := exit
	exit = func(code int) { exitCode = code; panic("exit") }
	t.Cleanup(func() { exit = oldExit })

	assert.PanicsWithValue(t, "exit", func() {
		Exit("fatal %s", "error")
	})

	assert.Equal(t, 1, exitCode)
	output := buf.String()
	assert.Contains(t, output, "ERROR")
	assert.Contains(t, output, "fatal error")
}

func TestPackageLevelInfoStartEnd(t *testing.T) {
	var buf bytes.Buffer
	l := NewWithWriter(false, 0, &buf)

	old := defaultLogger
	t.Cleanup(func() { defaultLogger = old })
	SetDefault(l)

	InfoStart("downloading")
	InfoEnd("done")

	output := buf.String()
	assert.Contains(t, output, "downloading")
	assert.Contains(t, output, "done")
	assert.Contains(t, output, "\n")
}
