package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		args     []string
		check    func(t *testing.T, cli *CLI)
		wantErr  bool
		errCheck func(t *testing.T, err error)
	}{
		// operations (short forms)
		{
			name: "update_short",
			args: []string{"-u"},
			check: func(t *testing.T, cli *CLI) {
				assert.True(t, cli.Update)
			},
		},
		{
			name: "list_short",
			args: []string{"-l"},
			check: func(t *testing.T, cli *CLI) {
				assert.True(t, cli.List)
			},
		},
		{
			name: "list_all_short",
			args: []string{"-a"},
			check: func(t *testing.T, cli *CLI) {
				assert.True(t, cli.ListAll)
			},
		},
		{
			name: "browse_short",
			args: []string{"-b"},
			check: func(t *testing.T, cli *CLI) {
				assert.True(t, cli.Browse)
			},
		},
		{
			name: "browse_long",
			args: []string{"--browse"},
			check: func(t *testing.T, cli *CLI) {
				assert.True(t, cli.Browse)
			},
		},
		{
			name: "search_short",
			args: []string{"-s", "ngi"},
			check: func(t *testing.T, cli *CLI) {
				assert.Equal(t, "ngi", cli.Search)
			},
		},
		{
			name: "info_short",
			args: []string{"-i"},
			check: func(t *testing.T, cli *CLI) {
				assert.True(t, cli.Info)
			},
		},
		{
			name: "render_short",
			args: []string{"-r", "file.md"},
			check: func(t *testing.T, cli *CLI) {
				assert.Equal(t, "file.md", cli.Render)
			},
		},

		// operations (long forms)
		{
			name: "update_long",
			args: []string{"--update"},
			check: func(t *testing.T, cli *CLI) {
				assert.True(t, cli.Update)
			},
		},
		{
			name: "list_long",
			args: []string{"--list"},
			check: func(t *testing.T, cli *CLI) {
				assert.True(t, cli.List)
			},
		},
		{
			name: "list_all_long",
			args: []string{"--list-all"},
			check: func(t *testing.T, cli *CLI) {
				assert.True(t, cli.ListAll)
			},
		},
		{
			name: "search_long",
			args: []string{"--search", "nginx"},
			check: func(t *testing.T, cli *CLI) {
				assert.Equal(t, "nginx", cli.Search)
			},
		},
		{
			name: "info_long",
			args: []string{"--info"},
			check: func(t *testing.T, cli *CLI) {
				assert.True(t, cli.Info)
			},
		},
		{
			name: "render_long",
			args: []string{"--render", "file.md"},
			check: func(t *testing.T, cli *CLI) {
				assert.Equal(t, "file.md", cli.Render)
			},
		},
		{
			name: "list_platforms",
			args: []string{"--list-platforms"},
			check: func(t *testing.T, cli *CLI) {
				assert.True(t, cli.ListPlatforms)
			},
		},
		{
			name: "list_languages",
			args: []string{"--list-languages"},
			check: func(t *testing.T, cli *CLI) {
				assert.True(t, cli.ListLanguages)
			},
		},
		{
			name: "clean_cache",
			args: []string{"--clean-cache"},
			check: func(t *testing.T, cli *CLI) {
				assert.True(t, cli.CleanCache)
			},
		},
		{
			name: "gen_config",
			args: []string{"--gen-config"},
			check: func(t *testing.T, cli *CLI) {
				assert.True(t, cli.GenConfig)
			},
		},
		{
			name: "config_path",
			args: []string{"--config-path"},
			check: func(t *testing.T, cli *CLI) {
				assert.True(t, cli.ConfigPath)
			},
		},

		// positional args
		{
			name: "single_page",
			args: []string{"tar"},
			check: func(t *testing.T, cli *CLI) {
				assert.Equal(t, []string{"tar"}, cli.Page)
			},
		},
		{
			name: "multiple_pages",
			args: []string{"tar", "git"},
			check: func(t *testing.T, cli *CLI) {
				assert.Equal(t, []string{"tar", "git"}, cli.Page)
			},
		},

		// options
		{
			name: "platform_short",
			args: []string{"-p", "linux", "-u"},
			check: func(t *testing.T, cli *CLI) {
				assert.Equal(t, "linux", cli.Platform)
			},
		},
		{
			name: "platform_long",
			args: []string{"--platform", "osx", "-u"},
			check: func(t *testing.T, cli *CLI) {
				assert.Equal(t, "osx", cli.Platform)
			},
		},
		{
			name: "language_single",
			args: []string{"-L", "en", "-u"},
			check: func(t *testing.T, cli *CLI) {
				assert.Equal(t, []string{"en"}, cli.Languages)
			},
		},
		{
			name: "language_repeat",
			args: []string{"-L", "en", "-L", "de", "-u"},
			check: func(t *testing.T, cli *CLI) {
				assert.Equal(t, []string{"en", "de"}, cli.Languages)
			},
		},
		{
			name: "language_comma",
			args: []string{"-L", "en,de", "-u"},
			check: func(t *testing.T, cli *CLI) {
				assert.Equal(t, []string{"en", "de"}, cli.Languages)
			},
		},
		{
			name: "language_long",
			args: []string{"--language", "fr", "-u"},
			check: func(t *testing.T, cli *CLI) {
				assert.Equal(t, []string{"fr"}, cli.Languages)
			},
		},
		{
			name: "offline_short",
			args: []string{"-o", "-u"},
			check: func(t *testing.T, cli *CLI) {
				assert.True(t, cli.Offline)
			},
		},
		{
			name: "offline_long",
			args: []string{"--offline", "-u"},
			check: func(t *testing.T, cli *CLI) {
				assert.True(t, cli.Offline)
			},
		},
		{
			name: "compact_short",
			args: []string{"-c", "-u"},
			check: func(t *testing.T, cli *CLI) {
				assert.True(t, cli.Compact)
			},
		},
		{
			name: "compact_long",
			args: []string{"--compact", "-u"},
			check: func(t *testing.T, cli *CLI) {
				assert.True(t, cli.Compact)
			},
		},
		{
			name: "no_compact",
			args: []string{"--no-compact", "-u"},
			check: func(t *testing.T, cli *CLI) {
				assert.True(t, cli.NoCompact)
			},
		},
		{
			name: "raw_short",
			args: []string{"-R", "-u"},
			check: func(t *testing.T, cli *CLI) {
				assert.True(t, cli.Raw)
			},
		},
		{
			name: "raw_long",
			args: []string{"--raw", "-u"},
			check: func(t *testing.T, cli *CLI) {
				assert.True(t, cli.Raw)
			},
		},
		{
			name: "no_raw",
			args: []string{"--no-raw", "-u"},
			check: func(t *testing.T, cli *CLI) {
				assert.True(t, cli.NoRaw)
			},
		},
		{
			name: "quiet_short",
			args: []string{"-q", "-u"},
			check: func(t *testing.T, cli *CLI) {
				assert.True(t, cli.Quiet)
			},
		},
		{
			name: "quiet_long",
			args: []string{"--quiet", "-u"},
			check: func(t *testing.T, cli *CLI) {
				assert.True(t, cli.Quiet)
			},
		},
		{
			name: "verbose_single",
			args: []string{"--verbose", "-u"},
			check: func(t *testing.T, cli *CLI) {
				assert.Equal(t, uint8(1), cli.Verbose)
			},
		},
		{
			name: "verbose_double",
			args: []string{"--verbose", "--verbose", "-u"},
			check: func(t *testing.T, cli *CLI) {
				assert.Equal(t, uint8(2), cli.Verbose)
			},
		},
		{
			name: "color_auto",
			args: []string{"--color", "auto", "-u"},
			check: func(t *testing.T, cli *CLI) {
				assert.Equal(t, "auto", cli.Color)
			},
		},
		{
			name: "color_always",
			args: []string{"--color", "always", "-u"},
			check: func(t *testing.T, cli *CLI) {
				assert.Equal(t, "always", cli.Color)
			},
		},
		{
			name: "color_never",
			args: []string{"--color", "never", "-u"},
			check: func(t *testing.T, cli *CLI) {
				assert.Equal(t, "never", cli.Color)
			},
		},
		{
			name: "color_default",
			args: []string{"-u"},
			check: func(t *testing.T, cli *CLI) {
				assert.Equal(t, "auto", cli.Color)
			},
		},
		{
			name: "config_path_option",
			args: []string{"--config", "/tmp/cfg", "-u"},
			check: func(t *testing.T, cli *CLI) {
				assert.Equal(t, "/tmp/cfg", cli.Config)
			},
		},
		{
			name: "edit",
			args: []string{"--edit", "-u"},
			check: func(t *testing.T, cli *CLI) {
				assert.True(t, cli.Edit)
			},
		},
		{
			name: "short_options",
			args: []string{"--short-options", "-u"},
			check: func(t *testing.T, cli *CLI) {
				assert.True(t, cli.ShortOptions)
			},
		},
		{
			name: "long_options",
			args: []string{"--long-options", "-u"},
			check: func(t *testing.T, cli *CLI) {
				assert.True(t, cli.LongOptions)
			},
		},

		// combined
		{
			name: "search_with_platform_and_language",
			args: []string{"-s", "ngi", "-p", "linux", "-L", "en"},
			check: func(t *testing.T, cli *CLI) {
				assert.Equal(t, "ngi", cli.Search)
				assert.Equal(t, "linux", cli.Platform)
				assert.Equal(t, []string{"en"}, cli.Languages)
			},
		},
		{
			name: "page_with_all_options",
			args: []string{"tar", "-p", "linux", "-L", "en", "-o", "-c", "-q"},
			check: func(t *testing.T, cli *CLI) {
				assert.Equal(t, []string{"tar"}, cli.Page)
				assert.Equal(t, "linux", cli.Platform)
				assert.Equal(t, []string{"en"}, cli.Languages)
				assert.True(t, cli.Offline)
				assert.True(t, cli.Compact)
				assert.True(t, cli.Quiet)
			},
		},

		// special operations
		{
			name: "version",
			args: []string{"-v"},
			check: func(t *testing.T, cli *CLI) {
				assert.True(t, cli.ShowVersion)
			},
		},
		{
			name: "version_long",
			args: []string{"--version"},
			check: func(t *testing.T, cli *CLI) {
				assert.True(t, cli.ShowVersion)
			},
		},
		{
			name: "help",
			args: []string{"-h"},
			check: func(t *testing.T, cli *CLI) {
				assert.True(t, cli.ShowHelp)
			},
		},
		{
			name: "help_long",
			args: []string{"--help"},
			check: func(t *testing.T, cli *CLI) {
				assert.True(t, cli.ShowHelp)
			},
		},
		{
			name:    "no_args_prints_help",
			args:    []string{},
			wantErr: false,
		},

		// error cases
		{
			name:    "browse_with_page_conflict",
			args:    []string{"-b", "tar"},
			wantErr: true,
			errCheck: func(t *testing.T, err error) {
				assert.ErrorContains(t, err, "cannot be used with")
			},
		},
		{
			name:    "two_operations",
			args:    []string{"-u", "-l"},
			wantErr: true,
			errCheck: func(t *testing.T, err error) {
				assert.ErrorContains(t, err, "cannot be used with")
			},
		},
		{
			name:    "three_operations",
			args:    []string{"-u", "-l", "-a"},
			wantErr: true,
			errCheck: func(t *testing.T, err error) {
				assert.ErrorContains(t, err, "cannot be used with")
			},
		},
		{
			name:    "invalid_color",
			args:    []string{"--color", "invalid", "-u"},
			wantErr: true,
			errCheck: func(t *testing.T, err error) {
				assert.ErrorContains(t, err, "invalid value")
			},
		},
		{
			name:    "unknown_flag",
			args:    []string{"--bogus"},
			wantErr: true,
			errCheck: func(t *testing.T, err error) {
				assert.ErrorContains(t, err, "unexpected argument")
			},
		},
		{
			name:    "unknown_short_flag",
			args:    []string{"-x"},
			wantErr: true,
			errCheck: func(t *testing.T, err error) {
				assert.ErrorContains(t, err, "unexpected argument")
			},
		},
		{
			name:    "unknown_flag_tip",
			args:    []string{"--searc"},
			wantErr: true,
			errCheck: func(t *testing.T, err error) {
				assert.ErrorContains(t, err, "unexpected argument")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli, err := parse(tt.args)
			if tt.wantErr {
				if tt.errCheck != nil {
					tt.errCheck(t, err)
				} else {
					assert.Error(t, err)
				}
				return
			}
			require.NoError(t, err)
			require.NotNil(t, cli)
			if tt.check != nil {
				tt.check(t, cli)
			}
		})
	}
}

func TestReorderFlags(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		args []string
		want []string
	}{
		{
			name: "empty",
			args: []string{},
			want: nil,
		},
		{
			name: "only_positional",
			args: []string{"page1", "page2"},
			want: []string{"page1", "page2"},
		},
		{
			name: "only_flags",
			args: []string{"-u"},
			want: []string{"-u"},
		},
		{
			name: "flags_before_positional",
			args: []string{"-u", "page1"},
			want: []string{"-u", "page1"},
		},
		{
			name: "flag_after_positional",
			args: []string{"page1", "-u"},
			want: []string{"-u", "page1"},
		},
		{
			name: "value_flag_after_positional",
			args: []string{"page1", "-p", "linux"},
			want: []string{"-p", "linux", "page1"},
		},
		{
			name: "long_flag_after_positional",
			args: []string{"page1", "--update"},
			want: []string{"--update", "page1"},
		},
		{
			name: "multiple_flags_after_positionals",
			args: []string{"page1", "-u", "-p", "linux", "-o"},
			want: []string{"-u", "-p", "linux", "-o", "page1"},
		},
		{
			name: "equals_syntax_stays_in_place",
			args: []string{"-p=linux", "page1"},
			want: []string{"-p=linux", "page1"},
		},
		{
			name: "mixed_positional_and_flags",
			args: []string{"page1", "-u", "page2", "-p", "linux"},
			want: []string{"-u", "-p", "linux", "page1", "page2"},
		},
		{
			name: "all_value_flags",
			args: []string{"page1", "-L", "en", "-s", "foo", "-r", "file.md", "--color", "always", "--config", "/tmp/cfg"},
			want: []string{"-L", "en", "-s", "foo", "-r", "file.md", "--color", "always", "--config", "/tmp/cfg", "page1"},
		},
		{
			name: "search_flag_after_positional",
			args: []string{"tar", "-s", "ngi"},
			want: []string{"-s", "ngi", "tar"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := reorderFlags(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}
