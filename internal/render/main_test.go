package render

// errorWriter is a test stub that always returns the configured write error.
type errorWriter struct {
	err error
}

// Write always returns the configured error.
func (w *errorWriter) Write(p []byte) (int, error) {
	return 0, w.err
}
