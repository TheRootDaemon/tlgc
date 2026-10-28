// Package exitcode defines the exit status codes used by tlgc.
package exitcode

const (
	// Success indicates that the command completed sucessfully.
	Success int = iota

	// InvalidArguments indicates that the command was invoked with invalid arguments.
	InvalidArguments

	// Default indicates a general IO, file, or other error.
	Default

	// CacheUpdateFailure indicates that updating the cache failed.
	CacheUpdateFailure

	// LintFailure indicates that the tldr page has lint error(s).
	LintFailure

	// ConfigParseFailure indicates that the configuration file could not be parsed.
	ConfigParseFailure

	// PageNotFound indicates that the requested page could not be found.
	PageNotFound
)
