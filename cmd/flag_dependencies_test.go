package cmd

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestModifierDependencies(t *testing.T) {
	t.Parallel()

	assert.NotEmpty(t, modifierDependencies)

	for _, dependency := range modifierDependencies {
		t.Run(dependency.flag, func(t *testing.T) {
			assert.True(t, strings.HasPrefix(dependency.flag, "--"))
			assert.NotEmpty(t, dependency.parents)

			for _, parent := range dependency.parents {
				assert.NotEmpty(t, parent)
			}

			assert.NotNil(t, dependency.present)
			assert.NotNil(t, dependency.valid)
		})
	}
}

func TestValidateFlagDependencies(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		cli       CLI
		wantErr   bool
		errString string
	}{
		// invalid: modifier present with no valid parent
		{
			name:      "output_without_parent",
			cli:       CLI{Output: "out.md"},
			wantErr:   true,
			errString: "flag --output requires --format",
		},
		{
			name:      "in_place_without_parent",
			cli:       CLI{InPlace: true},
			wantErr:   true,
			errString: "flag --in-place requires --format",
		},
		{
			name:      "tabular_without_parent",
			cli:       CLI{Tabular: true},
			wantErr:   true,
			errString: "flag --tabular requires --lint or --format",
		},
		{
			name:      "ignore_without_parent",
			cli:       CLI{Ignore: []string{"TLDR001"}},
			wantErr:   true,
			errString: "flag --ignore requires --lint or --format",
		},
		{
			name:      "platform_without_parent",
			cli:       CLI{Platform: "linux"},
			wantErr:   true,
			errString: "flag --platform requires a page, --browse, --list or --search",
		},
		{
			name:      "language_without_parent",
			cli:       CLI{Languages: []string{"en"}},
			wantErr:   true,
			errString: "flag --language requires a page, --browse, --search or --update",
		},
		{
			name:      "offline_without_parent",
			cli:       CLI{Offline: true},
			wantErr:   true,
			errString: "flag --offline requires a page or --browse",
		},
		{
			name:      "compact_without_parent",
			cli:       CLI{Compact: true},
			wantErr:   true,
			errString: "flag --compact requires a page or --render",
		},
		{
			name:      "no_compact_without_parent",
			cli:       CLI{NoCompact: true},
			wantErr:   true,
			errString: "flag --no-compact requires a page or --render",
		},
		{
			name:      "raw_without_parent",
			cli:       CLI{Raw: true},
			wantErr:   true,
			errString: "flag --raw requires a page or --render",
		},
		{
			name:      "no_raw_without_parent",
			cli:       CLI{NoRaw: true},
			wantErr:   true,
			errString: "flag --no-raw requires a page or --render",
		},
		{
			name:      "short_options_without_parent",
			cli:       CLI{ShortOptions: true},
			wantErr:   true,
			errString: "flag --short-options requires a page or --render",
		},
		{
			name:      "long_options_without_parent",
			cli:       CLI{LongOptions: true},
			wantErr:   true,
			errString: "flag --long-options requires a page or --render",
		},
		{
			name:      "edit_without_parent",
			cli:       CLI{Edit: true},
			wantErr:   true,
			errString: "flag --edit requires a page or --render",
		},
		{
			name:      "color_without_parent",
			cli:       CLI{Color: "always"},
			wantErr:   true,
			errString: "flag --color requires a page or --render",
		},
		{
			name: "color_auto_without_parent",
			cli:  CLI{Color: "auto"},
		},

		// valid: modifier with at least one active parent
		{
			name: "output_with_format",
			cli:  CLI{Output: "out.md", Format: true},
		},
		{
			name: "in_place_with_format",
			cli:  CLI{InPlace: true, Format: true},
		},
		{
			name: "tabular_with_lint",
			cli:  CLI{Tabular: true, Lint: true},
		},
		{
			name: "tabular_with_format",
			cli:  CLI{Tabular: true, Format: true},
		},
		{
			name: "ignore_with_lint",
			cli:  CLI{Ignore: []string{"TLDR001"}, Lint: true},
		},
		{
			name: "ignore_with_format",
			cli:  CLI{Ignore: []string{"TLDR001"}, Format: true},
		},
		{
			name: "platform_with_page",
			cli:  CLI{Platform: "linux", Page: []string{"tar"}},
		},
		{
			name: "platform_with_browse",
			cli:  CLI{Platform: "linux", Browse: true, Page: []string{"tar"}},
		},
		{
			name: "platform_with_list",
			cli:  CLI{Platform: "linux", List: true},
		},
		{
			name: "platform_with_search",
			cli:  CLI{Platform: "linux", Search: "ngi"},
		},
		{
			name: "language_with_page",
			cli:  CLI{Languages: []string{"en"}, Page: []string{"tar"}},
		},
		{
			name: "language_with_browse",
			cli:  CLI{Languages: []string{"en"}, Browse: true, Page: []string{"tar"}},
		},
		{
			name: "language_with_search",
			cli:  CLI{Languages: []string{"en"}, Search: "ngi"},
		},
		{
			name: "language_with_update",
			cli:  CLI{Languages: []string{"en"}, Update: true},
		},
		{
			name: "offline_with_page",
			cli:  CLI{Offline: true, Page: []string{"tar"}},
		},
		{
			name: "offline_with_browse",
			cli:  CLI{Offline: true, Browse: true, Page: []string{"tar"}},
		},
		{
			name: "compact_with_page",
			cli:  CLI{Compact: true, Page: []string{"tar"}},
		},
		{
			name: "compact_with_render",
			cli:  CLI{Compact: true, Render: "file.md"},
		},
		{
			name: "no_compact_with_page",
			cli:  CLI{NoCompact: true, Page: []string{"tar"}},
		},
		{
			name: "no_compact_with_render",
			cli:  CLI{NoCompact: true, Render: "file.md"},
		},
		{
			name: "raw_with_page",
			cli:  CLI{Raw: true, Page: []string{"tar"}},
		},
		{
			name: "raw_with_render",
			cli:  CLI{Raw: true, Render: "file.md"},
		},
		{
			name: "no_raw_with_page",
			cli:  CLI{NoRaw: true, Page: []string{"tar"}},
		},
		{
			name: "no_raw_with_render",
			cli:  CLI{NoRaw: true, Render: "file.md"},
		},
		{
			name: "short_options_with_page",
			cli:  CLI{ShortOptions: true, Page: []string{"tar"}},
		},
		{
			name: "short_options_with_render",
			cli:  CLI{ShortOptions: true, Render: "file.md"},
		},
		{
			name: "long_options_with_page",
			cli:  CLI{LongOptions: true, Page: []string{"tar"}},
		},
		{
			name: "long_options_with_render",
			cli:  CLI{LongOptions: true, Render: "file.md"},
		},
		{
			name: "edit_with_page",
			cli:  CLI{Edit: true, Page: []string{"tar"}},
		},
		{
			name: "edit_with_render",
			cli:  CLI{Edit: true, Render: "file.md"},
		},
		{
			name: "color_with_page",
			cli:  CLI{Color: "always", Page: []string{"tar"}},
		},
		{
			name: "color_with_render",
			cli:  CLI{Color: "never", Render: "file.md"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli := tt.cli
			if cli.Color == "" {
				cli.Color = "auto"
			}
			err := validateFlagDependencies(&cli)
			if tt.wantErr {
				assert.ErrorContains(t, err, tt.errString)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPageLookup(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		cli  CLI
		want bool
	}{
		{
			name: "no_page",
			cli:  CLI{},
			want: false,
		},
		{
			name: "page",
			cli:  CLI{Page: []string{"tar"}},
			want: true,
		},
		{
			name: "multiple_pages",
			cli:  CLI{Page: []string{"tar", "git"}},
			want: true,
		},
		{
			name: "page_with_browse",
			cli:  CLI{Page: []string{"tar"}, Browse: true},
			want: false,
		},
		{
			name: "page_with_lint",
			cli:  CLI{Page: []string{"file.md"}, Lint: true},
			want: false,
		},
		{
			name: "page_with_format",
			cli:  CLI{Page: []string{"file.md"}, Format: true},
			want: false,
		},
		{
			name: "browse_without_page",
			cli:  CLI{Browse: true},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.cli.pageLookup())
		})
	}
}

func TestFormatParents(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		items []string
		want  string
	}{
		{
			name:  "single",
			items: []string{"--format"},
			want:  "--format",
		},
		{
			name:  "two",
			items: []string{"--lint", "--format"},
			want:  "--lint or --format",
		},
		{
			name:  "three",
			items: []string{"a page", "--browse", "--render"},
			want:  "a page, --browse or --render",
		},
		{
			name:  "four",
			items: []string{"a page", "--browse", "--list", "--search"},
			want:  "a page, --browse, --list or --search",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, formatParents(tt.items))
		})
	}
}
