package upstream

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type closeTracker struct {
	io.ReadCloser
	closeCalled bool
}

func (t *closeTracker) Close() error {
	t.closeCalled = true
	return t.ReadCloser.Close()
}

func TestNewRequest(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		client *Client
		ctx    context.Context
		url    string
		check  func(t *testing.T, req *http.Request, err error)
	}{
		{
			name:   "valid url with custom user agent",
			client: &Client{userAgent: "my-agent/42"},
			ctx:    context.Background(),
			url:    "https://example.com/file",
			check: func(t *testing.T, req *http.Request, err error) {
				require.NoError(t, err)
				require.NotNil(t, req)
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "https://example.com/file", req.URL.String())
				assert.Equal(t, "my-agent/42", req.Header.Get("User-Agent"))
			},
		},
		{
			name:   "valid url with default user agent",
			client: &Client{userAgent: defaultUserAgent},
			ctx:    context.Background(),
			url:    "https://example.com/file",
			check: func(t *testing.T, req *http.Request, err error) {
				require.NoError(t, err)
				require.NotNil(t, req)
				assert.Equal(t, defaultUserAgent, req.Header.Get("User-Agent"))
			},
		},
		{
			name:   "empty url",
			client: &Client{userAgent: defaultUserAgent},
			ctx:    context.Background(),
			url:    "",
			check: func(t *testing.T, req *http.Request, err error) {
				require.NoError(t, err)
				require.NotNil(t, req)
				assert.Equal(t, "", req.URL.String())
			},
		},
		{
			name:   "invalid url scheme",
			client: &Client{userAgent: defaultUserAgent},
			ctx:    context.Background(),
			url:    "://invalid",
			check: func(t *testing.T, req *http.Request, err error) {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "missing protocol scheme")
				assert.Nil(t, req)
			},
		},
		{
			name:   "pre-cancelled context still creates request",
			client: &Client{userAgent: defaultUserAgent},
			ctx: func() context.Context {
				cctx, cancel := context.WithCancel(context.Background())
				cancel()
				return cctx
			}(),
			url: "https://example.com/file",
			check: func(t *testing.T, req *http.Request, err error) {
				require.NoError(t, err)
				require.NotNil(t, req)
				assert.Error(t, req.Context().Err())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req, err := tt.client.newRequest(tt.ctx, tt.url)
			tt.check(t, req, err)
		})
	}
}

func TestSend(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		handler http.HandlerFunc
		ctx     context.Context
		url     string
		check   func(t *testing.T, resp *http.Response, err error)
	}{
		{
			name: "200 ok",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				_, _ = w.Write([]byte("ok"))
			},
			check: func(t *testing.T, resp *http.Response, err error) {
				require.NoError(t, err)
				require.NotNil(t, resp)
				assert.Equal(t, 200, resp.StatusCode)
				_ = resp.Body.Close()
			},
		},
		{
			name: "404 not found",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(404)
			},
			check: func(t *testing.T, resp *http.Response, err error) {
				require.NoError(t, err)
				require.NotNil(t, resp)
				assert.Equal(t, 404, resp.StatusCode)
				_ = resp.Body.Close()
			},
		},
		{
			name: "cancelled context returns error",
			ctx: func() context.Context {
				cctx, cancel := context.WithCancel(context.Background())
				cancel()
				return cctx
			}(),
			url: "http://127.0.0.1:1",
			check: func(t *testing.T, resp *http.Response, err error) {
				require.Error(t, err)
				assert.Nil(t, resp)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var ts *httptest.Server
			if tt.handler != nil {
				ts = httptest.NewServer(tt.handler)
				defer ts.Close()
			}

			reqURL := tt.url
			if ts != nil {
				reqURL = ts.URL
			}

			c := &Client{
				client:    http.DefaultClient,
				userAgent: defaultUserAgent,
			}

			ctx := tt.ctx
			if ctx == nil {
				ctx = context.Background()
			}

			req, err := c.newRequest(ctx, reqURL)
			require.NoError(t, err)
			require.NotNil(t, req)

			resp, err := c.send(req)
			tt.check(t, resp, err)
		})
	}
}

func TestValidateResponse(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		code    int
		wantErr string
	}{
		{name: "200 ok", code: 200},
		{name: "201 created", code: 201},
		{name: "299 non-standard success", code: 299},
		{name: "300 multiple choices", code: 300, wantErr: "unexpected status: 300"},
		{name: "301 moved permanently", code: 301, wantErr: "unexpected status: 301"},
		{name: "400 bad request", code: 400, wantErr: "unexpected status: 400"},
		{name: "401 unauthorized", code: 401, wantErr: "unexpected status: 401"},
		{name: "403 forbidden", code: 403, wantErr: "unexpected status: 403"},
		{name: "404 not found", code: 404, wantErr: "unexpected status: 404"},
		{name: "500 internal server error", code: 500, wantErr: "unexpected status: 500"},
		{name: "199 informational", code: 199, wantErr: "unexpected status: 199"},
		{name: "100 continue", code: 100, wantErr: "unexpected status: 100"},
	}

	c := &Client{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			resp := &http.Response{
				StatusCode: tt.code,
				Body:       io.NopCloser(strings.NewReader("")),
			}

			err := c.validateResponse(resp)
			if tt.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestWrapBody(t *testing.T) {
	t.Parallel()

	t.Run("no wrapping when maxBodySize=0 and progress=nil", func(t *testing.T) {
		t.Parallel()

		body := io.NopCloser(strings.NewReader("hello"))
		c := &Client{}

		wrapped := c.wrapBody(body, 5)
		assert.Equal(t, body, wrapped)

		data, err := io.ReadAll(wrapped)
		require.NoError(t, err)
		assert.Equal(t, "hello", string(data))

		_ = wrapped.Close()
	})

	t.Run("limitedBody wraps when maxBodySize>0 and progress=nil", func(t *testing.T) {
		t.Parallel()

		body := io.NopCloser(strings.NewReader("hello world"))
		c := &Client{maxBodySize: 5}

		wrapped := c.wrapBody(body, 11)
		assert.NotEqual(t, body, wrapped)

		data, err := io.ReadAll(wrapped)
		require.NoError(t, err)
		assert.Equal(t, 6, len(data))
		assert.Equal(t, "hello ", string(data))

		_ = wrapped.Close()
	})

	t.Run("progressReader wraps when progress!=nil and maxBodySize=0", func(t *testing.T) {
		t.Parallel()

		var mu sync.Mutex
		var done, total int64

		body := io.NopCloser(strings.NewReader("hello"))
		c := &Client{
			progress: func(d, t int64) {
				mu.Lock()
				done = d
				total = t
				mu.Unlock()
			},
		}

		wrapped := c.wrapBody(body, 5)
		assert.NotEqual(t, body, wrapped)

		data, err := io.ReadAll(wrapped)
		require.NoError(t, err)
		assert.Equal(t, "hello", string(data))

		mu.Lock()
		assert.Equal(t, int64(5), done)
		assert.Equal(t, int64(5), total)
		mu.Unlock()

		_ = wrapped.Close()
	})

	t.Run("both wrappers applied when both conditions are true", func(t *testing.T) {
		t.Parallel()

		var mu sync.Mutex
		var done, total int64

		body := io.NopCloser(strings.NewReader("hello world this is a long body"))
		c := &Client{
			maxBodySize: 10,
			progress: func(d, t int64) {
				mu.Lock()
				done = d
				total = t
				mu.Unlock()
			},
		}

		wrapped := c.wrapBody(body, 21)
		assert.NotEqual(t, body, wrapped)

		data, err := io.ReadAll(wrapped)
		require.NoError(t, err)
		assert.Equal(t, 11, len(data))
		assert.Equal(t, "hello world", string(data))

		mu.Lock()
		assert.Equal(t, int64(11), done)
		assert.Equal(t, int64(21), total)
		mu.Unlock()

		_ = wrapped.Close()
	})

	t.Run("empty body with limitedBody", func(t *testing.T) {
		t.Parallel()

		body := io.NopCloser(strings.NewReader(""))
		c := &Client{maxBodySize: 100}

		wrapped := c.wrapBody(body, 0)
		assert.NotEqual(t, body, wrapped)

		data, err := io.ReadAll(wrapped)
		require.NoError(t, err)
		assert.Equal(t, 0, len(data))

		_ = wrapped.Close()
	})

	t.Run("close propagates through limitedBody", func(t *testing.T) {
		t.Parallel()

		body := &closeTracker{
			ReadCloser: io.NopCloser(strings.NewReader("hello world")),
		}
		c := &Client{maxBodySize: 10}

		wrapped := c.wrapBody(body, 11)
		assert.NotEqual(t, body, wrapped)

		err := wrapped.Close()
		require.NoError(t, err)
		assert.True(t, body.closeCalled)
	})

	t.Run("close propagates through progressReader", func(t *testing.T) {
		t.Parallel()

		body := &closeTracker{
			ReadCloser: io.NopCloser(strings.NewReader("hello")),
		}
		c := &Client{
			progress: func(d, t int64) {},
		}

		wrapped := c.wrapBody(body, 5)
		assert.NotEqual(t, body, wrapped)

		err := wrapped.Close()
		require.NoError(t, err)
		assert.True(t, body.closeCalled)
	})

	t.Run("close propagates through both wrappers", func(t *testing.T) {
		t.Parallel()

		body := &closeTracker{
			ReadCloser: io.NopCloser(strings.NewReader("hello world this is long")),
		}
		c := &Client{
			maxBodySize: 10,
			progress:    func(d, t int64) {},
		}

		wrapped := c.wrapBody(body, 25)
		assert.NotEqual(t, body, wrapped)

		err := wrapped.Close()
		require.NoError(t, err)
		assert.True(t, body.closeCalled)
	})
}

func TestExecute(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		ctx     context.Context
		url     string
		opts    []Option
		handler http.HandlerFunc
		check   func(t *testing.T, resp *http.Response, err error)
	}{
		{
			name: "200 with body",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				_, _ = w.Write([]byte("hello"))
			},
			check: func(t *testing.T, resp *http.Response, err error) {
				require.NoError(t, err)
				require.NotNil(t, resp)
				assert.Equal(t, 200, resp.StatusCode)
				body, _ := io.ReadAll(resp.Body)
				_ = resp.Body.Close()
				assert.Equal(t, "hello", string(body))
			},
		},
		{
			name: "200 empty body",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
			},
			check: func(t *testing.T, resp *http.Response, err error) {
				require.NoError(t, err)
				require.NotNil(t, resp)
				body, _ := io.ReadAll(resp.Body)
				_ = resp.Body.Close()
				assert.Equal(t, "", string(body))
			},
		},
		{
			name: "404 not found",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(404)
			},
			check: func(t *testing.T, resp *http.Response, err error) {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "unexpected status: 404")
				assert.Nil(t, resp)
			},
		},
		{
			name: "500 internal server error",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(500)
			},
			check: func(t *testing.T, resp *http.Response, err error) {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "unexpected status: 500")
				assert.Nil(t, resp)
			},
		},
		{
			name: "custom user agent is sent",
			opts: []Option{WithUserAgent("my-agent/42")},
			handler: func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte(r.Header.Get("User-Agent")))
			},
			check: func(t *testing.T, resp *http.Response, err error) {
				require.NoError(t, err)
				require.NotNil(t, resp)
				body, _ := io.ReadAll(resp.Body)
				_ = resp.Body.Close()
				assert.Equal(t, "my-agent/42", string(body))
			},
		},
		{
			name: "body limited by maxBodySize",
			opts: []Option{WithMaxBodySize(10)},
			handler: func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte("hello world this is a long body"))
			},
			check: func(t *testing.T, resp *http.Response, err error) {
				require.NoError(t, err)
				require.NotNil(t, resp)
				body, _ := io.ReadAll(resp.Body)
				_ = resp.Body.Close()
				assert.Equal(t, 11, len(body))
				assert.Equal(t, "hello world", string(body))
			},
		},
		{
			name: "cancelled context returns error",
			ctx: func() context.Context {
				cctx, cancel := context.WithCancel(context.Background())
				cancel()
				return cctx
			}(),
			url: "https://example.com/file",
			check: func(t *testing.T, resp *http.Response, err error) {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "context canceled")
				assert.Nil(t, resp)
			},
		},
		{
			name: "invalid url",
			url:  "://invalid",
			check: func(t *testing.T, resp *http.Response, err error) {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "missing protocol scheme")
				assert.Nil(t, resp)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var ts *httptest.Server
			if tt.handler != nil {
				ts = httptest.NewServer(tt.handler)
				defer ts.Close()
			}

			reqURL := tt.url
			if ts != nil {
				reqURL = ts.URL
			}

			c := New(tt.opts...)

			ctx := tt.ctx
			if ctx == nil {
				ctx = context.Background()
			}

			resp, err := c.execute(ctx, reqURL)
			tt.check(t, resp, err)
		})
	}
}

func TestProgressCallback_WithContentLength(t *testing.T) {
	t.Parallel()

	var mu sync.Mutex
	var capturedDone, capturedTotal int64

	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Length", "5")
				_, _ = w.Write([]byte("hello"))
			},
		),
	)
	defer ts.Close()

	c := &Client{
		client:    ts.Client(),
		userAgent: defaultUserAgent,
		progress: func(done, total int64) {
			mu.Lock()
			capturedDone = done
			capturedTotal = total
			mu.Unlock()
		},
	}

	resp, err := c.execute(context.Background(), ts.URL)
	require.NoError(t, err)
	require.NotNil(t, resp)

	body, err := io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	require.NoError(t, err)
	require.Equal(t, "hello", string(body))

	mu.Lock()
	assert.Equal(t, int64(5), capturedDone)
	assert.Equal(t, int64(5), capturedTotal)
	mu.Unlock()
}

func TestProgressCallback_WithUnknownContentLength(t *testing.T) {
	t.Parallel()

	var mu sync.Mutex
	var capturedDone, capturedTotal int64

	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				flusher := w.(http.Flusher)
				_, _ = w.Write([]byte("ab"))
				flusher.Flush()
				_, _ = w.Write([]byte("cd"))
				flusher.Flush()
			},
		),
	)
	defer ts.Close()

	c := &Client{
		client:    ts.Client(),
		userAgent: defaultUserAgent,
		progress: func(done, total int64) {
			mu.Lock()
			capturedDone = done
			capturedTotal = total
			mu.Unlock()
		},
	}

	resp, err := c.execute(context.Background(), ts.URL)
	require.NoError(t, err)
	require.NotNil(t, resp)

	body, _ := io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	assert.Equal(t, "abcd", string(body))

	mu.Lock()
	assert.Equal(t, int64(4), capturedDone)
	assert.Equal(t, int64(-1), capturedTotal)
	mu.Unlock()
}
