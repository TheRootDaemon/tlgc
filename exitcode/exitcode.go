// Package exitcode defines the exit status codes used by tlgc.
package exitcode

const (
	// Success indicates that the command completed sucessfully.
	Success int = iota

	// InvalidArguments indicates that the command was invoked with invalid arguments.
	InvalidArguments

	// CacheUpdateFailure indicates that updating the cache failed.
	CacheUpdateFailure

	// LintFailure indicates that the tldr page has lint error(s).
	LintFailure
)
