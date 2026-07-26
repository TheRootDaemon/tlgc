package cmd

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimilarFlag(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		flagSet []string
		query   string
		want    string
	}{
		{
			name:    "prefix_priority",
			flagSet: []string{"search", "update"},
			query:   "sea",
			want:    "search",
		},
		{
			name:    "fuzzy_fallback",
			flagSet: []string{"search", "update"},
			query:   "searh",
			want:    "search",
		},
		{
			name:    "no_match",
			flagSet: []string{"update", "list"},
			query:   "xyz",
			want:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := flag.NewFlagSet("test", flag.ContinueOnError)
			for _, name := range tt.flagSet {
				fs.Bool(name, false, "")
			}
			got := similarFlag(fs, tt.query)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPrefixFlag(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		flagSet []string
		query   string
		want    string
	}{
		{
			name:    "exact_prefix",
			flagSet: []string{"update", "list", "list-all"},
			query:   "lis",
			want:    "list",
		},
		{
			name:    "shortest_prefix_wins",
			flagSet: []string{"list", "list-all", "list-platforms"},
			query:   "list",
			want:    "list",
		},
		{
			name:    "no_match",
			flagSet: []string{"update", "list"},
			query:   "bogus",
			want:    "",
		},
		{
			name:    "skips_short_flags",
			flagSet: []string{"u", "q"},
			query:   "u",
			want:    "",
		},
		{
			name:    "single_match",
			flagSet: []string{"search"},
			query:   "sea",
			want:    "search",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := flag.NewFlagSet("test", flag.ContinueOnError)
			for _, name := range tt.flagSet {
				fs.Bool(name, false, "")
			}
			got := prefixFlag(fs, tt.query)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFuzzyFlag(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		flagSet []string
		query   string
		want    string
	}{
		{
			name:    "distance_1",
			flagSet: []string{"search", "update"},
			query:   "searh",
			want:    "search",
		},
		{
			name:    "distance_2",
			flagSet: []string{"verbose", "offline"},
			query:   "verbos",
			want:    "verbose",
		},
		{
			name:    "no_match_too_far",
			flagSet: []string{"update", "list"},
			query:   "xyz",
			want:    "",
		},
		{
			name:    "skips_short_flags",
			flagSet: []string{"u", "update"},
			query:   "updaet",
			want:    "update",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := flag.NewFlagSet("test", flag.ContinueOnError)
			for _, name := range tt.flagSet {
				fs.Bool(name, false, "")
			}
			got := fuzzyFlag(fs, tt.query)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestEditDistance(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		a, b string
		want int
	}{
		{name: "identical", a: "abc", b: "abc", want: 0},
		{name: "empty_both", a: "", b: "", want: 0},
		{name: "empty_a", a: "", b: "abc", want: 3},
		{name: "empty_b", a: "abc", b: "", want: 3},
		{name: "substitution", a: "abc", b: "axc", want: 1},
		{name: "insertion", a: "ac", b: "abc", want: 1},
		{name: "deletion", a: "abc", b: "ac", want: 1},
		{name: "full_mismatch", a: "abc", b: "xyz", want: 3},
		{name: "typo", a: "searc", b: "search", want: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := editDistance(tt.a, tt.b)
			assert.Equal(t, tt.want, got)
		})
	}
}
