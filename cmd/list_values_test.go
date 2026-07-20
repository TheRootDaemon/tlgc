package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCountValueNilPointer(t *testing.T) {
	t.Parallel()

	v := &countValue{}
	assert.Equal(t, "0", v.String())
}

func TestCountValue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		setCount int
		want     string
		wantBool bool
	}{
		{
			name:     "zero_value",
			want:     "0",
			wantBool: true,
		},
		{
			name:     "single_increment",
			setCount: 1,
			want:     "1",
			wantBool: true,
		},
		{
			name:     "double_increment",
			setCount: 2,
			want:     "2",
			wantBool: true,
		},
		{
			name:     "triple_increment",
			setCount: 3,
			want:     "3",
			wantBool: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var count uint8
			v := &countValue{count: &count}

			for range tt.setCount {
				require.NoError(t, v.Set(""))
			}

			assert.Equal(t, tt.want, v.String())
			assert.Equal(t, tt.wantBool, v.IsBoolFlag())
		})
	}
}

func TestStringListValueNilPointer(t *testing.T) {
	t.Parallel()

	v := &stringListValue{}
	assert.Equal(t, "", v.String())
}

func TestStringListValue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		inputs []string
		want   string
		wantS  []string
	}{
		{
			name:   "empty",
			inputs: nil,
			want:   "",
			wantS:  nil,
		},
		{
			name:   "single_value",
			inputs: []string{"en"},
			want:   "en",
			wantS:  []string{"en"},
		},
		{
			name:   "comma_separated",
			inputs: []string{"de,pl"},
			want:   "de,pl",
			wantS:  []string{"de", "pl"},
		},
		{
			name:   "whitespace_trimmed",
			inputs: []string{" de , pl "},
			want:   "de,pl",
			wantS:  []string{"de", "pl"},
		},
		{
			name:   "empty_parts_skipped",
			inputs: []string{",en,,de,"},
			want:   "en,de",
			wantS:  []string{"en", "de"},
		},
		{
			name:   "multiple_calls",
			inputs: []string{"en", "de"},
			want:   "en,de",
			wantS:  []string{"en", "de"},
		},
		{
			name:   "mixed_comma_and_repeated",
			inputs: []string{"en,de", "fr"},
			want:   "en,de,fr",
			wantS:  []string{"en", "de", "fr"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var values []string
			v := &stringListValue{values: &values}

			for _, input := range tt.inputs {
				require.NoError(t, v.Set(input))
			}

			assert.Equal(t, tt.want, v.String())
			assert.Equal(t, tt.wantS, values)
		})
	}
}
