package upstream

import (
	"io"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProgressReader_Read(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		data          string
		total         int64
		readSize      int
		wantData      string
		wantDone      int64
		wantCallCount int
	}{
		{
			name:          "single_read",
			data:          "hello",
			total:         5,
			readSize:      64,
			wantData:      "hello",
			wantDone:      5,
			wantCallCount: 2,
		},
		{
			name:          "multiple_reads",
			data:          "abcdef",
			total:         6,
			readSize:      2,
			wantData:      "abcdef",
			wantDone:      6,
			wantCallCount: 4,
		},
		{
			name:          "unknown_length",
			data:          "xyz",
			total:         -1,
			readSize:      64,
			wantData:      "xyz",
			wantDone:      3,
			wantCallCount: 2,
		},
		{
			name:          "empty_data",
			data:          "",
			total:         0,
			readSize:      64,
			wantData:      "",
			wantDone:      0,
			wantCallCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var mu sync.Mutex
			var calls []struct{ done, total int64 }

			pr := &progressReader{
				reader: io.NopCloser(strings.NewReader(tt.data)),
				total:  tt.total,
				progressCallback: func(done, total int64) {
					mu.Lock()
					calls = append(calls, struct{ done, total int64 }{done, total})
					mu.Unlock()
				},
			}
			defer func() {
				_ = pr.Close()
			}()

			buf := make([]byte, tt.readSize)
			var result strings.Builder
			for {
				n, err := pr.Read(buf)
				result.Write(buf[:n])
				if err != nil {
					break
				}
			}

			assert.Equal(t, tt.wantData, result.String())
			assert.Equal(t, tt.wantDone, pr.done)
			assert.Equal(t, tt.wantCallCount, len(calls))

			mu.Lock()
			if len(calls) > 0 {
				assert.Equal(t, tt.total, calls[len(calls)-1].total)
			}
			mu.Unlock()
		})
	}
}

func TestProgressReader_NilFn(t *testing.T) {
	t.Parallel()

	pr := &progressReader{
		reader: io.NopCloser(strings.NewReader("hello")),
		total:  5,
	}
	defer func() {
		_ = pr.Close()
	}()

	buf := make([]byte, 64)
	require.NotPanics(t, func() {
		n, err := pr.Read(buf)
		require.NoError(t, err)
		assert.Equal(t, "hello", string(buf[:n]))
		assert.Equal(t, int64(5), pr.done)
	})
}

func TestLimitedBody_Interface(t *testing.T) {
	t.Parallel()

	var _ io.ReadCloser = (*limitedBody)(nil)
}

func TestLimitedBody_Read(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		content  string
		limit    int64
		wantData string
	}{
		{
			name:     "under_limit",
			content:  "hello",
			limit:    100,
			wantData: "hello",
		},
		{
			name:     "exactly_limit_plus_one",
			content:  "hello world",
			limit:    10,
			wantData: "hello world",
		},
		{
			name:     "empty_body",
			content:  "",
			limit:    100,
			wantData: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			lb := &limitedBody{
				Reader: io.LimitReader(strings.NewReader(tt.content), tt.limit+1),
				closer: io.NopCloser(strings.NewReader("")),
				limit:  tt.limit,
			}
			defer func() {
				_ = lb.Close()
			}()

			data, err := io.ReadAll(lb)
			require.NoError(t, err)
			assert.Equal(t, tt.wantData, string(data))
		})
	}
}
