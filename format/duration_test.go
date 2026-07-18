package format

import (
	"math"
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

func TestValidateDurationOverflow(t *testing.T) {
	t.Parallel()

	const maxValid = math.MaxInt64 / int64(time.Hour)

	tests := []struct {
		name    string
		hours   uint64
		want    time.Duration
		wantErr string
	}{
		{name: "zero_hours", hours: 0, want: 0},
		{name: "one_hour", hours: 1, want: time.Hour},
		{name: "max_valid", hours: uint64(maxValid), want: time.Duration(maxValid) * time.Hour},
		{name: "overflow", hours: uint64(maxValid) + 1, wantErr: "overflows"},
		{name: "max_uint64", hours: math.MaxUint64, wantErr: "overflows"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidateDurationOverflow(tt.hours)
			if tt.wantErr != "" {
				assert.ErrorContains(t, err, tt.wantErr)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
