// Package exitcode defines the exit status codes used by tlgc.
package exitcode

const (
	// Success indicates that the command completed sucessfully.
	Success = 0

	// InvalidArguments indicates that the command was invoked with invalid arguments.
	InvalidArguments = 1

	// CacheUpdateFailure indicates that updating the cache failed.
	CacheUpdateFailure = 2
)
