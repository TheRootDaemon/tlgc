package upstream

import (
	"io"
)

type progressReader struct {
	reader io.ReadCloser
	total  int64
	done   int64
	fn     func(done, total int64)
}

func (r *progressReader) Read(p []byte) (int, error) {
	n, err := r.reader.Read(p)
	r.done += int64(n)
	if r.fn != nil {
		r.fn(r.done, r.total)
	}

	return n, err
}

func (r *progressReader) Close() error {
	return r.reader.Close()
}

type limitedBody struct {
	io.Reader
	closer io.Closer
	limit  int64
}

func (b *limitedBody) Close() error {
	return b.closer.Close()
}
