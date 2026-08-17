package app

import (
	"testing"

	"github.com/TheRootDaemon/tlgc/cmd"
	"github.com/stretchr/testify/assert"
)

func TestShouldAutoUpdate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		cli  *cmd.CLI
		want bool
	}{
		{
			name: "default_invocation",
			cli:  &cmd.CLI{},
			want: true,
		},
		{
			name: "offline_suppresses",
			cli:  &cmd.CLI{Offline: true},
			want: false,
		},
		{
			name: "explicit_update_suppresses",
			cli:  &cmd.CLI{Update: true},
			want: false,
		},
		{
			name: "clean_cache_suppresses",
			cli:  &cmd.CLI{CleanCache: true},
			want: false,
		},
		{
			name: "offline_beats_update",
			cli: &cmd.CLI{
				Offline: true,
				Update:  true,
			}, want: false,
		},
		{
			name: "offline_beats_clean",
			cli: &cmd.CLI{
				Offline:    true,
				CleanCache: true,
			}, want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, shouldAutoUpdate(tt.cli))
		})
	}
}
