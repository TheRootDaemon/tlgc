package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		cli     CLI
		wantErr bool
	}{
		{
			name: "single_operation_update",
			cli:  CLI{Color: "auto", Update: true},
		},
		{
			name: "single_operation_list",
			cli:  CLI{Color: "auto", List: true},
		},
		{
			name: "single_operation_search",
			cli:  CLI{Color: "auto", Search: "ngi"},
		},
		{
			name:    "no_operations",
			cli:     CLI{Color: "auto"},
			wantErr: false,
		},
		{
			name:    "no_operations_with_modifiers",
			cli:     CLI{Color: "auto", HasArgs: true, Compact: true, Edit: true},
			wantErr: true,
		},
		{
			name:    "invalid_color",
			cli:     CLI{Color: "invalid", Update: true},
			wantErr: true,
		},
		{
			name:    "two_operations",
			cli:     CLI{Color: "auto", Update: true, List: true},
			wantErr: true,
		},
		{
			name:    "three_operations",
			cli:     CLI{Color: "auto", Update: true, List: true, ListAll: true},
			wantErr: true,
		},
		{
			name: "browse_with_page",
			cli:  CLI{Color: "auto", Browse: true, Page: []string{"tar"}},
		},
		{
			name:    "page_and_update",
			cli:     CLI{Color: "auto", Page: []string{"tar"}, Update: true},
			wantErr: true,
		},
		{
			name:    "browse_and_update",
			cli:     CLI{Color: "auto", Browse: true, Update: true},
			wantErr: true,
		},
		{
			name:    "browse_without_page",
			cli:     CLI{Color: "auto", Browse: true},
			wantErr: true,
		},
		{
			name:    "lint_without_path",
			cli:     CLI{Color: "auto", Lint: true},
			wantErr: true,
		},
		{
			name: "lint_with_path",
			cli:  CLI{Color: "auto", Lint: true, Page: []string{"pages/"}},
		},
		{
			name:    "format_without_path",
			cli:     CLI{Color: "auto", Format: true},
			wantErr: true,
		},
		{
			name: "format_with_path",
			cli:  CLI{Color: "auto", Format: true, Page: []string{"file.md"}},
		},
		{
			name:    "output_without_format",
			cli:     CLI{Color: "auto", Output: "out.md", Page: []string{"file.md"}},
			wantErr: true,
		},
		{
			name: "output_with_format",
			cli:  CLI{Color: "auto", Format: true, Output: "out.md", Page: []string{"file.md"}},
		},
		{
			name:    "in_place_without_format",
			cli:     CLI{Color: "auto", InPlace: true, Page: []string{"file.md"}},
			wantErr: true,
		},
		{
			name: "in_place_with_format",
			cli:  CLI{Color: "auto", Format: true, InPlace: true, Page: []string{"file.md"}},
		},
		{
			name: "tabular_with_lint",
			cli:  CLI{Color: "auto", Lint: true, Tabular: true, Page: []string{"file.md"}},
		},
		{
			name: "tabular_with_format",
			cli:  CLI{Color: "auto", Format: true, Tabular: true, Page: []string{"file.md"}},
		},
		{
			name:    "tabular_alone",
			cli:     CLI{Color: "auto", HasArgs: true, Tabular: true},
			wantErr: true,
		},
		{
			name: "ignore_with_lint",
			cli:  CLI{Color: "auto", Lint: true, Ignore: []string{"TLDR001"}, Page: []string{"file.md"}},
		},
		{
			name: "ignore_with_format",
			cli:  CLI{Color: "auto", Format: true, Ignore: []string{"TLDR001"}, Page: []string{"file.md"}},
		},
		{
			name:    "ignore_alone",
			cli:     CLI{Color: "auto", HasArgs: true, Ignore: []string{"TLDR001"}},
			wantErr: true,
		},
		{
			name: "platform_with_page",
			cli:  CLI{Color: "auto", Platform: "linux", Page: []string{"tar"}},
		},
		{
			name: "platform_with_browse",
			cli:  CLI{Color: "auto", Platform: "linux", Browse: true, Page: []string{"tar"}},
		},
		{
			name: "platform_with_list",
			cli:  CLI{Color: "auto", Platform: "linux", List: true},
		},
		{
			name: "platform_with_search",
			cli:  CLI{Color: "auto", Platform: "linux", Search: "ngi"},
		},
		{
			name:    "platform_with_update",
			cli:     CLI{Color: "auto", Platform: "linux", Update: true},
			wantErr: true,
		},
		{
			name: "language_with_update",
			cli:  CLI{Color: "auto", Languages: []string{"en"}, Update: true},
		},
		{
			name: "language_with_search",
			cli:  CLI{Color: "auto", Languages: []string{"en"}, Search: "ngi"},
		},
		{
			name:    "language_with_list",
			cli:     CLI{Color: "auto", Languages: []string{"en"}, List: true},
			wantErr: true,
		},
		{
			name: "offline_with_browse",
			cli:  CLI{Color: "auto", Offline: true, Browse: true, Page: []string{"tar"}},
		},
		{
			name:    "offline_with_update",
			cli:     CLI{Color: "auto", Offline: true, Update: true},
			wantErr: true,
		},
		{
			name: "compact_with_render",
			cli:  CLI{Color: "auto", Compact: true, Render: "file.md"},
		},
		{
			name:    "compact_with_search",
			cli:     CLI{Color: "auto", Compact: true, Search: "ngi"},
			wantErr: true,
		},
		{
			name: "edit_with_render",
			cli:  CLI{Color: "auto", Edit: true, Render: "file.md"},
		},
		{
			name: "color_with_page",
			cli:  CLI{Color: "always", Page: []string{"tar"}},
		},
		{
			name:    "color_with_update",
			cli:     CLI{Color: "always", Update: true},
			wantErr: true,
		},
		{
			name:    "valid_color_auto",
			cli:     CLI{Color: "auto", Update: true},
			wantErr: false,
		},
		{
			name:    "valid_color_always",
			cli:     CLI{Color: "always", Page: []string{"tar"}},
			wantErr: false,
		},
		{
			name:    "valid_color_never",
			cli:     CLI{Color: "never", Page: []string{"tar"}},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate(&tt.cli)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateColor(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		color     string
		wantErr   bool
		errString string
	}{
		{
			name:  "auto",
			color: "auto",
		},
		{
			name:  "always",
			color: "always",
		},
		{
			name:  "never",
			color: "never",
		},
		{
			name:      "invalid",
			color:     "invalid",
			wantErr:   true,
			errString: "invalid value for --color",
		},
		{
			name:      "empty",
			color:     "",
			wantErr:   true,
			errString: "invalid value for --color",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateColor(&CLI{Color: tt.color})
			if tt.wantErr {
				assert.ErrorContains(t, err, tt.errString)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateOperations(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		cli       CLI
		wantErr   bool
		errString string
	}{
		{
			name: "one_operation",
			cli:  CLI{HasArgs: true, Update: true},
		},
		{
			name:      "no_operations_with_args",
			cli:       CLI{HasArgs: true},
			wantErr:   true,
			errString: "no operation specified",
		},
		{
			name:      "two_operations",
			cli:       CLI{HasArgs: true, Update: true, List: true},
			wantErr:   true,
			errString: "cannot be used with",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateOperations(&tt.cli)
			if tt.wantErr {
				assert.ErrorContains(t, err, tt.errString)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateOperationArguments(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		cli       CLI
		wantErr   bool
		errString string
	}{
		{
			name: "browse_with_page",
			cli:  CLI{Browse: true, Page: []string{"tar"}},
		},
		{
			name:      "browse_without_page",
			cli:       CLI{Browse: true},
			wantErr:   true,
			errString: "flag --browse requires a page",
		},
		{
			name: "lint_with_path",
			cli:  CLI{Lint: true, Page: []string{"pages/"}},
		},
		{
			name:      "lint_without_path",
			cli:       CLI{Lint: true},
			wantErr:   true,
			errString: "flag --lint requires a file or directory",
		},
		{
			name: "format_with_path",
			cli:  CLI{Format: true, Page: []string{"file.md"}},
		},
		{
			name:      "format_without_path",
			cli:       CLI{Format: true},
			wantErr:   true,
			errString: "flag --format requires a file or directory",
		},
		{
			name: "output_with_format",
			cli:  CLI{Format: true, Output: "out.md", Page: []string{"file.md"}},
		},
		{
			name:      "output_without_format",
			cli:       CLI{Output: "out.md", Page: []string{"file.md"}},
			wantErr:   true,
			errString: "flag --output requires --format",
		},
		{
			name: "no_operation",
			cli:  CLI{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateOperationArguments(&tt.cli)
			if tt.wantErr {
				assert.ErrorContains(t, err, tt.errString)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestOperationCount(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		cli  CLI
		want int
	}{
		{
			name: "none",
			cli:  CLI{},
			want: 0,
		},
		{
			name: "page",
			cli:  CLI{Page: []string{"tar"}},
			want: 1,
		},
		{
			name: "update",
			cli:  CLI{Update: true},
			want: 1,
		},
		{
			name: "list",
			cli:  CLI{List: true},
			want: 1,
		},
		{
			name: "list_all",
			cli:  CLI{ListAll: true},
			want: 1,
		},
		{
			name: "search",
			cli:  CLI{Search: "ngi"},
			want: 1,
		},
		{
			name: "browse",
			cli:  CLI{Browse: true, Page: []string{"tar"}},
			want: 1,
		},
		{
			name: "list_platforms",
			cli:  CLI{ListPlatforms: true},
			want: 1,
		},
		{
			name: "list_languages",
			cli:  CLI{ListLanguages: true},
			want: 1,
		},
		{
			name: "info",
			cli:  CLI{Info: true},
			want: 1,
		},
		{
			name: "render",
			cli:  CLI{Render: "file.md"},
			want: 1,
		},
		{
			name: "lint",
			cli:  CLI{Lint: true, Page: []string{"file.md"}},
			want: 1,
		},
		{
			name: "format",
			cli:  CLI{Format: true, Page: []string{"file.md"}},
			want: 1,
		},
		{
			name: "lint_and_format",
			cli:  CLI{Lint: true, Format: true, Page: []string{"file.md"}},
			want: 2,
		},
		{
			name: "clean_cache",
			cli:  CLI{CleanCache: true},
			want: 1,
		},
		{
			name: "gen_config",
			cli:  CLI{GenConfig: true},
			want: 1,
		},
		{
			name: "config_path",
			cli:  CLI{ConfigPath: true},
			want: 1,
		},
		{
			name: "two_operations",
			cli:  CLI{Update: true, List: true},
			want: 2,
		},
		{
			name: "all_operations",
			cli: CLI{
				Page:          []string{"tar"},
				Update:        true,
				List:          true,
				ListAll:       true,
				Search:        "ngi",
				Browse:        true,
				ListPlatforms: true,
				ListLanguages: true,
				Info:          true,
				Render:        "file.md",
				CleanCache:    true,
				GenConfig:     true,
				ConfigPath:    true,
			},
			want: 12,
		},
		{
			name: "all_lint_operations",
			cli: CLI{
				Lint:    true,
				Format:  true,
				Output:  "out.md",
				InPlace: true,
				Tabular: true,
				Ignore:  []string{"TLDR001"},
				Page:    []string{"file.md"},
			},
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.cli.operationCount()
			assert.Equal(t, tt.want, got)
		})
	}
}
