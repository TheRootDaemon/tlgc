package app

import (
	"testing"

	"github.com/TheRootDaemon/tlgc/cmd"
	"github.com/TheRootDaemon/tlgc/internal/config"
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
			want: false,
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
		{
			name: "page_lookup",
			cli:  &cmd.CLI{Page: []string{"tar"}},
			want: true,
		},
		{
			name: "browse",
			cli:  &cmd.CLI{Browse: true, Page: []string{"tar"}},
			want: true,
		},
		{
			name: "search",
			cli:  &cmd.CLI{Search: "ngi"},
			want: true,
		},
		{
			name: "list",
			cli:  &cmd.CLI{List: true},
			want: true,
		},
		{
			name: "help_suppresses",
			cli:  &cmd.CLI{ShowHelp: true},
			want: false,
		},
		{
			name: "version_suppresses",
			cli:  &cmd.CLI{ShowVersion: true},
			want: false,
		},
		{
			name: "info",
			cli:  &cmd.CLI{Info: true},
			want: true,
		},
		{
			name: "list_all_suppresses",
			cli:  &cmd.CLI{ListAll: true},
			want: false,
		},
		{
			name: "list_platforms_suppresses",
			cli:  &cmd.CLI{ListPlatforms: true},
			want: false,
		},
		{
			name: "list_languages_suppresses",
			cli:  &cmd.CLI{ListLanguages: true},
			want: false,
		},
		{
			name: "render_suppresses",
			cli:  &cmd.CLI{Render: "file.md"},
			want: false,
		},
		{
			name: "lint_suppresses",
			cli:  &cmd.CLI{Lint: true, Page: []string{"pages/"}},
			want: false,
		},
		{
			name: "format_suppresses",
			cli:  &cmd.CLI{Format: true, Page: []string{"file.md"}},
			want: false,
		},
		{
			name: "gen_config_suppresses",
			cli:  &cmd.CLI{GenConfig: true},
			want: false,
		},
		{
			name: "config_path_suppresses",
			cli:  &cmd.CLI{ConfigPath: true},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, shouldAutoUpdate(tt.cli))
		})
	}
}

func TestShouldAutoUpdateDisabledByConfig(t *testing.T) {
	cfg := config.Default()
	cfg.Cache.AutoUpdate = false

	config.SetForTesting(&cfg)
	defer config.ResetForTesting()

	assert.False(t, shouldAutoUpdate(&cmd.CLI{Search: "ngi"}))
	assert.False(t, shouldAutoUpdate(&cmd.CLI{Page: []string{"tar"}}))
}
