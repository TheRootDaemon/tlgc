package logger

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrace(t *testing.T) {
	var buf bytes.Buffer

	l := &Logger{
		level:  LevelTrace,
		writer: &buf,
	}

	l.Trace("hello %s", "world")

	output := buf.String()

	assert.Contains(t, output, "trace:")
	assert.Contains(t, output, "hello world")
	assert.True(t, strings.HasSuffix(output, "\n"))
}

func TestTraceDisabled(t *testing.T) {
	var buf bytes.Buffer

	l := &Logger{
		level:  LevelDebug,
		writer: &buf,
	}

	l.Trace("hidden")

	assert.Empty(t, buf.String())
}

func TestDebug(t *testing.T) {
	var buf bytes.Buffer

	l := &Logger{
		level:  LevelDebug,
		writer: &buf,
	}

	l.Debug("hello %s", "world")

	output := buf.String()

	assert.Contains(t, output, "debug:")
	assert.Contains(t, output, "hello world")
	assert.True(t, strings.HasSuffix(output, "\n"))
}

func TestDebugDisabled(t *testing.T) {
	var buf bytes.Buffer

	l := &Logger{
		level:  LevelInfo,
		writer: &buf,
	}

	l.Debug("hidden")

	assert.Empty(t, buf.String())
}

func TestInfo(t *testing.T) {
	var buf bytes.Buffer

	l := &Logger{
		level:  LevelInfo,
		writer: &buf,
	}

	l.Info("hello %s", "world")

	output := buf.String()

	assert.Contains(t, output, "info:")
	assert.Contains(t, output, "hello world")
	assert.True(t, strings.HasSuffix(output, "\n"))
}

func TestInfoDisabled(t *testing.T) {
	var buf bytes.Buffer

	l := &Logger{
		level:  LevelWarn,
		writer: &buf,
	}

	l.Info("hidden")

	assert.Empty(t, buf.String())
}

func TestWarn(t *testing.T) {
	var buf bytes.Buffer

	l := &Logger{
		level:  LevelWarn,
		writer: &buf,
	}

	l.Warn("hello %s", "world")

	output := buf.String()

	assert.Contains(t, output, "warning:")
	assert.Contains(t, output, "hello world")
	assert.True(t, strings.HasSuffix(output, "\n"))
}

func TestWarnDisabled(t *testing.T) {
	var buf bytes.Buffer

	l := &Logger{
		level:  LevelError,
		writer: &buf,
	}

	l.Warn("hidden")

	assert.Empty(t, buf.String())
}

func TestError(t *testing.T) {
	var buf bytes.Buffer

	l := &Logger{
		level:  LevelError,
		writer: &buf,
	}

	l.Error("hello %s", "world")

	output := buf.String()

	assert.Contains(t, output, "error:")
	assert.Contains(t, output, "hello world")
	assert.True(t, strings.HasSuffix(output, "\n"))
}

func TestLog(t *testing.T) {
	var buf bytes.Buffer

	l := &Logger{
		level:  LevelInfo,
		writer: &buf,
	}

	l.Log(LevelInfo, "hello %s", "world")

	output := buf.String()

	assert.Contains(t, output, "info:")
	assert.Contains(t, output, "hello world")
	assert.True(t, strings.HasSuffix(output, "\n"))
}

func TestLogDisabled(t *testing.T) {
	var buf bytes.Buffer

	l := &Logger{
		level:  LevelWarn,
		writer: &buf,
	}

	l.Log(LevelInfo, "hidden")

	assert.Empty(t, buf.String())
}

func TestInfoStartEnd(t *testing.T) {
	tests := []struct {
		name    string
		quiet   bool
		verbose uint8
	}{
		{name: "info mode", quiet: false, verbose: 0},
		{name: "debug mode", quiet: false, verbose: 1},
		{name: "quiet mode", quiet: true, verbose: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			l := NewWithWriter(tt.quiet, tt.verbose, &buf)
			l.InfoStart("downloading")
			l.InfoEnd("done")
			output := buf.String()

			switch {
			case tt.quiet:
				assert.Empty(t, output, "quiet mode should suppress")
			case tt.verbose >= 1:
				assert.Contains(t, output, "downloading")
				assert.Contains(t, output, "\n")
				assert.NotContains(t, output, "\ndone", "end should be no-op in debug mode")
			default:
				assert.Contains(t, output, "downloadingdone", "should be one line")
				assert.Contains(t, output, "\n", "should end with newline")
			}
		})
	}
}

func TestInfoStartFallback(t *testing.T) {
	var buf bytes.Buffer
	l := NewWithWriter(false, 1, &buf)
	l.InfoStart("fallback")

	assert.Contains(t, buf.String(), "fallback")
	assert.Contains(t, buf.String(), "\n")
}

func TestInfoEndNoopWithoutStart(t *testing.T) {
	var buf bytes.Buffer
	l := NewWithWriter(false, 1, &buf)
	l.InfoEnd("orphan")

	assert.Empty(t, buf.String())
}
