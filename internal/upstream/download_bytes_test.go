package upstream

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDownloadBytes(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				_, _ = w.Write([]byte("hello world"))
			},
		),
	)

	defer ts.Close()

	c := New()
	ctx := context.Background()

	data, err := c.DownloadBytes(ctx, ts.URL, "")
	require.NoError(t, err)
	assert.Equal(t, "hello world", string(data))
}
