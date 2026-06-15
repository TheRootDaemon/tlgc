package duration

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDurationFmt(t *testing.T) {
	tests := []struct {
		name  string
		input time.Duration
		want  string
	}{
		{name: "seconds only", input: 1 * time.Second, want: "1s"},
		{name: "exact minute", input: 1 * time.Minute, want: "1min"},
		{name: "minute and seconds", input: 1*time.Minute + 1*time.Second, want: "1min, 1s"},
		{name: "exact hour", input: 1 * time.Hour, want: "1h"},
		{name: "hour with seconds truncated", input: 1*time.Hour + 1*time.Second, want: "1h"},
		{name: "hour and minute", input: 1*time.Hour + 1*time.Minute, want: "1h, 1min"},
		{name: "hour minute second", input: 1*time.Hour + 1*time.Minute + 1*time.Second, want: "1h, 1min"},
		{name: "exact day", input: 24 * time.Hour, want: "1d"},
		{name: "day with seconds truncated", input: 24*time.Hour + 1*time.Second, want: "1d"},
		{name: "day and hour", input: 24*time.Hour + 1*time.Hour, want: "1d, 1h"},
		{name: "day hour second truncated", input: 24*time.Hour + 1*time.Hour + 1*time.Second, want: "1d, 1h"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DurationFmt(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}
