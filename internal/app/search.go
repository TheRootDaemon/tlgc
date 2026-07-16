package app

import (
	"fmt"
	"strings"

	"github.com/TheRootDaemon/tlgc/cmd"
	"github.com/TheRootDaemon/tlgc/internal/cache"
	"github.com/TheRootDaemon/tlgc/logger"
	"github.com/TheRootDaemon/tlgc/termcolor"
)

// searchPages searches cached pages for the given query.
// Returns 0 on success, 1 on error.
func (a *App) searchPages(cli *cmd.CLI) int {
	c := cache.New()
	p := a.resolvePlatform(cli.Platform)
	languages := a.resolveLanguages(cli.Languages)

	results, err := c.Search(cli.Search, p, languages)
	if err != nil {
		logger.Error("search failed: %v", err)
		return 1
	}

	if err := a.printSearchResults(results, cli.Search); err != nil {
		logger.Error("%w", err)
		return 1
	}

	return 0
}

// printSearchResults renders the search results as a table,
// including the header and all matching rows.
func (a *App) printSearchResults(
	results []cache.SearchResult,
	query string,
) error {
	langW, platW, pageW := columnWidths(results)

	if err := a.printSearchHeader(langW, platW, pageW); err != nil {
		return err
	}

	return a.printSearchRows(
		results,
		query,
		langW,
		platW,
		pageW,
	)
}

// printSearchHeader writes the table header
// using the provided column widths.
func (a *App) printSearchHeader(langW, platW, pageW int) error {
	header := termcolor.Fprintf(
		"bold",
		"%-*s %-*s %-*s",
		langW,
		"Language",
		platW,
		"Platform",
		pageW,
		"Page",
	)

	_, err := fmt.Fprintln(a.Stdout, header)
	return err
}

// printSearchRows writes each search result as a table row,
// highlighting occurrences of query in the page name.
func (a *App) printSearchRows(
	results []cache.SearchResult,
	query string,
	langW, platW, pageW int,
) error {
	for _, r := range results {
		page := highlightQuery(r.Page, query)
		padding := max(pageW-len(r.Page), 0)

		_, err := fmt.Fprintf(
			a.Stdout,
			"%-*s %-*s %s%s\n",
			langW,
			r.Language,
			platW,
			r.Platform,
			page,
			strings.Repeat(" ", padding),
		)
		if err != nil {
			return err
		}
	}

	return nil
}

// columnWidths returns the widths required to display the
// Language, Platform, and Page columns without truncation.
func columnWidths(results []cache.SearchResult) (int, int, int) {
	langW := len("Language")
	platW := len("Platform")
	pageW := len("Page")

	for _, r := range results {
		langW = max(langW, len(r.Language))
		platW = max(platW, len(r.Platform))
		pageW = max(pageW, len(r.Page))
	}

	return langW, platW, pageW
}

// highlightQuery returns text
// with every case-insensitive occurrence of query
// wrapped in ANSI color codes.
func highlightQuery(text, query string) string {
	if query == "" {
		return text
	}

	lower := strings.ToLower(text)
	queryLower := strings.ToLower(query)

	var b strings.Builder
	b.Grow(len(text) + len(query)*10)

	remaining := text
	rest := lower
	for {
		idx := strings.Index(rest, queryLower)
		if idx < 0 {
			b.WriteString(remaining)
			break
		}

		b.WriteString(remaining[:idx])
		b.WriteString(
			termcolor.Sprint(
				"blue",
				remaining[idx:idx+len(query)],
			),
		)

		remaining = remaining[idx+len(query):]
		rest = rest[idx+len(query):]
	}

	return b.String()
}
