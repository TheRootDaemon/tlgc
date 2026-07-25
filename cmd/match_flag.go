package cmd

import (
	"flag"
	"strings"
)

// similarFlag finds the closest defined long flag to name, within threshold.
func similarFlag(fs *flag.FlagSet, name string) string {
	if prefix := prefixFlag(fs, name); prefix != "" {
		return prefix
	}
	return fuzzyFlag(fs, name)
}

// prefixFlag returns the shortest defined long flag that starts with name.
func prefixFlag(fs *flag.FlagSet, name string) string {
	var best string

	fs.VisitAll(func(f *flag.Flag) {
		if len(f.Name) <= 1 {
			return
		}

		if strings.HasPrefix(f.Name, name) &&
			(best == "" || len(f.Name) < len(best)) {
			best = f.Name
		}
	})

	return best
}

// fuzzyFlag returns the closest defined long flag to name by Levenshtein distance.
func fuzzyFlag(fs *flag.FlagSet, name string) string {
	var best string
	bestDist := 3

	fs.VisitAll(func(f *flag.Flag) {
		if len(f.Name) <= 1 {
			return
		}

		d := editDistance(name, f.Name)
		if d < bestDist {
			bestDist = d
			best = f.Name
		}
	})

	return best
}

// editDistance returns the Levenshtein distance between a and b.
func editDistance(a, b string) int {
	rows, cols := len(a)+1, len(b)+1
	distances := make([][]int, rows)

	for i := range distances {
		distances[i] = make([]int, cols)
		distances[i][0] = i
	}

	for j := range cols {
		distances[0][j] = j
	}

	for i := 1; i < rows; i++ {
		for j := 1; j < cols; j++ {
			cost := 0
			if a[i-1] != b[j-1] {
				cost = 1
			}
			distances[i][j] = min(
				distances[i-1][j]+1,
				distances[i][j-1]+1,
				distances[i-1][j-1]+cost,
			)
		}
	}

	return distances[rows-1][cols-1]
}
