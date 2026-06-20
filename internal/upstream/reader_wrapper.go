package upstream

import (
	"io"
)

// progressReader wraps an io.ReadCloser
// and reports cumulative read progress
// after each Read call.
type progressReader struct {
	reader           io.ReadCloser
	total            int64
	done             int64
	progressCallback func(done, total int64)
}

// Read reads from the underlying reader,
// updates the cumulative byte count,
// and invokes the progress callback if configured.
func (r *progressReader) Read(p []byte) (int, error) {
	n, err := r.reader.Read(p)
	r.done += int64(n)
	if r.progressCallback != nil {
		r.progressCallback(r.done, r.total)
	}

	return n, err
}

// Close closes the underlying reader.
func (r *progressReader) Close() error {
	return r.reader.Close()
}

// limitedBody wraps a reader with an associated closer,
// allowing read limits to be enforced
// while preserving Close behavior.
type limitedBody struct {
	io.Reader
	closer io.Closer
	limit  int64
}

// Close closes the underlying reader.
func (b *limitedBody) Close() error {
	return b.closer.Close()
}
