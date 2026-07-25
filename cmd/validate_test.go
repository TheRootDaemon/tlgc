package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		cli         CLI
		wantErr     bool
		errContains string
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
			name:        "invalid_color",
			cli:         CLI{Color: "invalid", Update: true},
			wantErr:     true,
			errContains: "invalid value",
		},
		{
			name:        "two_operations",
			cli:         CLI{Color: "auto", Update: true, List: true},
			wantErr:     true,
			errContains: "cannot be used with",
		},
		{
			name:        "three_operations",
			cli:         CLI{Color: "auto", Update: true, List: true, ListAll: true},
			wantErr:     true,
			errContains: "cannot be used with",
		},
		{
			name:        "page_and_update",
			cli:         CLI{Color: "auto", Page: []string{"tar"}, Update: true},
			wantErr:     true,
			errContains: "cannot be used with",
		},
		{
			name:    "valid_color_auto",
			cli:     CLI{Color: "auto", Update: true},
			wantErr: false,
		},
		{
			name:    "valid_color_always",
			cli:     CLI{Color: "always", Update: true},
			wantErr: false,
		},
		{
			name:    "valid_color_never",
			cli:     CLI{Color: "never", Update: true},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate(&tt.cli)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.ErrorContains(t, err, tt.errContains)
				}
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.cli.operationCount()
			assert.Equal(t, tt.want, got)
		})
	}
}
