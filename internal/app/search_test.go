package app

import (
	"testing"

	"github.com/TheRootDaemon/tlgc/internal/cache"
	"github.com/stretchr/testify/assert"
)

func TestColumnWidths(t *testing.T) {
	t.Parallel()

	languageHeader := len("Language")
	platformHeader := len("Platform")
	pageHeader := len("Page")

	tests := []struct {
		name     string
		results  []cache.SearchResult
		wantLang int
		wantPlat int
		wantPage int
	}{
		{
			name:     "empty_results_uses_header_widths",
			results:  []cache.SearchResult{},
			wantLang: languageHeader,
			wantPlat: platformHeader,
			wantPage: pageHeader,
		},
		{
			name: "header_wins_over_short_values",
			results: []cache.SearchResult{
				{Language: "en", Platform: "common", Page: "tar"},
			},
			wantLang: languageHeader,
			wantPlat: platformHeader,
			wantPage: pageHeader,
		},
		{
			name: "long_page_widens_page_column",
			results: []cache.SearchResult{
				{Language: "en", Platform: "common", Page: "very-long-page-name"},
			},
			wantLang: languageHeader,
			wantPlat: platformHeader,
			wantPage: len("very-long-page-name"),
		},
		{
			name: "long_language_widens_language_column",
			results: []cache.SearchResult{
				{Language: "pt_BR", Platform: "common", Page: "tar"},
			},
			wantLang: languageHeader,
			wantPlat: platformHeader,
			wantPage: pageHeader,
		},
		{
			name: "long_platform_widens_platform_column",
			results: []cache.SearchResult{
				{Language: "en", Platform: "android", Page: "tar"},
			},
			wantLang: languageHeader,
			wantPlat: platformHeader,
			wantPage: pageHeader,
		},
		{
			name: "picks_max_across_all_rows",
			results: []cache.SearchResult{
				{Language: "en", Platform: "common", Page: "apt"},
				{Language: "fr_FR", Platform: "android", Page: "very-long-page-name"},
			},
			wantLang: languageHeader,
			wantPlat: platformHeader,
			wantPage: len("very-long-page-name"),
		},
		{
			name: "multiple_languages_picks_widest",
			results: []cache.SearchResult{
				{Language: "en", Platform: "common", Page: "tar"},
				{Language: "de", Platform: "common", Page: "git"},
				{Language: "pt_BR", Platform: "common", Page: "ls"},
			},
			wantLang: languageHeader,
			wantPlat: platformHeader,
			wantPage: pageHeader,
		},
		{
			name: "value_exceeding_header_widens_column",
			results: []cache.SearchResult{
				{Language: "en", Platform: "common", Page: "tar"},
				{Language: "en_AUSTRALIA", Platform: "common", Page: "git"},
			},
			wantLang: len("en_AUSTRALIA"),
			wantPlat: platformHeader,
			wantPage: pageHeader,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLang, gotPlat, gotPage := columnWidths(tt.results)
			assert.Equal(t, tt.wantLang, gotLang)
			assert.Equal(t, tt.wantPlat, gotPlat)
			assert.Equal(t, tt.wantPage, gotPage)
		})
	}
}

func TestHighlightQuery(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		text  string
		query string
	}{
		{
			name:  "empty_query",
			text:  "nginx",
			query: "",
		},
		{
			name:  "no_match",
			text:  "nginx",
			query: "xyz",
		},
		{
			name:  "single_match_at_start",
			text:  "nginx",
			query: "ngi",
		},
		{
			name:  "single_match_at_end",
			text:  "nginx",
			query: "inx",
		},
		{
			name:  "single_match_in_middle",
			text:  "git-commit",
			query: "com",
		},
		{
			name:  "multiple_matches",
			text:  "nginx nginx",
			query: "ng",
		},
		{
			name:  "full_match",
			text:  "git",
			query: "git",
		},
		{
			name:  "case_insensitive_match",
			text:  "Nginx",
			query: "ngi",
		},
		{
			name:  "mixed_case_query",
			text:  "nginx",
			query: "NgI",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := highlightQuery(tt.text, tt.query)
			assert.Equal(t, tt.text, stripANSI(got), "visible text must be preserved")
		})
	}
}

// stripANSI removes ANSI escape sequences from s.
func stripANSI(s string) string {
	var out []byte
	inEscape := false
	for i := 0; i < len(s); i++ {
		switch {
		case s[i] == '\x1b':
			inEscape = true
		case inEscape && s[i] == 'm':
			inEscape = false
		case !inEscape:
			out = append(out, s[i])
		}
	}
	return string(out)
}
