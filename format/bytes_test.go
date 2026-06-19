package format

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBytesFmt(t *testing.T) {
	tests := []struct {
		name  string
		input int64
		want  string
	}{
		{name: "zero", input: 0, want: "0 B"},
		{name: "one byte", input: 1, want: "1 B"},
		{name: "just below 1 KiB", input: 1023, want: "1023 B"},
		{name: "exactly 1 KiB", input: 1024, want: "1.0 KiB"},
		{name: "just above 1 KiB", input: 1025, want: "1.0 KiB"},
		{name: "1.5 KiB", input: 1536, want: "1.5 KiB"},
		{name: "exactly 1 MiB", input: 1048576, want: "1.0 MiB"},
		{name: "1.5 MiB", input: 1572864, want: "1.5 MiB"},
		{name: "exactly 1 GiB", input: 1073741824, want: "1.0 GiB"},
		{name: "1.5 GiB", input: 1610612736, want: "1.5 GiB"},
		{name: "exactly 1 TiB", input: int64(1099511627776), want: "1.0 TiB"},
		{name: "1.5 TiB", input: int64(1649267441664), want: "1.5 TiB"},
		{name: "exactly 1 PiB", input: int64(1125899906842624), want: "1.0 PiB"},
		{name: "exactly 1 EiB", input: int64(1152921504606846976), want: "1.0 EiB"},
		{name: "max int64", input: math.MaxInt64, want: "8.0 EiB"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BytesFmt(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}
