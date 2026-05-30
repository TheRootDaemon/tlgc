package logger

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLevelConstants(t *testing.T) {
	assert.Less(t, LevelTrace, LevelDebug)
	assert.Less(t, LevelDebug, LevelInfo)
	assert.Less(t, LevelInfo, LevelWarn)
	assert.Less(t, LevelWarn, LevelError)
}

func TestLevelLabelsExist(t *testing.T) {
	levels := []Level{LevelTrace, LevelDebug, LevelInfo, LevelWarn, LevelError}
	for _, lvl := range levels {
		_, okColor := levelColors[lvl]
		assert.True(t, okColor, "missing color for level %d", lvl)
		_, okLabel := levelLabels[lvl]
		assert.True(t, okLabel, "missing label for level %d", lvl)
	}
}

func TestNew(t *testing.T) {
	l := New(false, 0)
	assert.Equal(t, LevelInfo, l.level)
	assert.NotNil(t, l.writer)
}

func TestNewWithWriter(t *testing.T) {
	tests := []struct {
		name    string
		quiet   bool
		verbose uint8
		want    Level
	}{
		{name: "default", quiet: false, verbose: 0, want: LevelInfo},
		{name: "quiet", quiet: true, verbose: 0, want: LevelError},
		{name: "quiet with verbose", quiet: true, verbose: 1, want: LevelError},
		{name: "verbose 1", quiet: false, verbose: 1, want: LevelDebug},
		{name: "verbose 2", quiet: false, verbose: 2, want: LevelTrace},
		{name: "verbose 3", quiet: false, verbose: 3, want: LevelTrace},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewWithWriter(tt.quiet, tt.verbose, new(bytes.Buffer))
			assert.Equal(t, tt.want, l.level)
		})
	}
}

func TestEnabled(t *testing.T) {
	tests := []struct {
		name  string
		level Level
		arg   Level
		want  bool
	}{
		{name: "trace when error", level: LevelError, arg: LevelTrace, want: false},
		{name: "debug when error", level: LevelError, arg: LevelDebug, want: false},
		{name: "info when error", level: LevelError, arg: LevelInfo, want: false},
		{name: "warn when error", level: LevelError, arg: LevelWarn, want: false},
		{name: "error when error", level: LevelError, arg: LevelError, want: true},
		{name: "trace when info", level: LevelInfo, arg: LevelTrace, want: false},
		{name: "debug when info", level: LevelInfo, arg: LevelDebug, want: false},
		{name: "info when info", level: LevelInfo, arg: LevelInfo, want: true},
		{name: "warn when info", level: LevelInfo, arg: LevelWarn, want: true},
		{name: "error when info", level: LevelInfo, arg: LevelError, want: true},
		{name: "trace when trace", level: LevelTrace, arg: LevelTrace, want: true},
		{name: "debug when trace", level: LevelTrace, arg: LevelDebug, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Logger{level: tt.level}
			assert.Equal(t, tt.want, l.Enabled(tt.arg))
		})
	}
}
