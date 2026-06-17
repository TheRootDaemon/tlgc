package upstream

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Run("defaults", func(t *testing.T) {
		c := New()

		assert.NotNil(t, c)
		assert.Equal(t, defaultTimeout, c.timeout)
		assert.Equal(t, defaultConnectionTimeout, c.connectionTimeout)
		assert.Equal(t, defaultUserAgent, c.userAgent)
		assert.Equal(t, int64(defaultMaxBodySize), c.maxBodySize)
		assert.Nil(t, c.progress)
		assert.NotNil(t, c.client)
		assert.Equal(t, defaultTimeout, c.client.Timeout)
	})
}

func TestWithOptions(t *testing.T) {
	customClient := &http.Client{Timeout: 123 * time.Second}

	tests := []struct {
		name  string
		opts  []Option
		check func(t *testing.T, c *Client)
	}{
		{
			name: "WithTimeout",
			opts: []Option{WithTimeout(5 * time.Second)},
			check: func(t *testing.T, c *Client) {
				assert.Equal(t, 5*time.Second, c.timeout)
				assert.Equal(t, 5*time.Second, c.client.Timeout)
			},
		},
		{
			name: "WithConnectionTimeout",
			opts: []Option{WithConnectionTimeout(3 * time.Second)},
			check: func(t *testing.T, c *Client) {
				assert.Equal(t, 3*time.Second, c.connectionTimeout)
			},
		},
		{
			name: "WithUserAgent",
			opts: []Option{WithUserAgent("test-agent/1.0")},
			check: func(t *testing.T, c *Client) {
				assert.Equal(t, "test-agent/1.0", c.userAgent)
			},
		},
		{
			name: "WithMaxBodySize",
			opts: []Option{WithMaxBodySize(1024)},
			check: func(t *testing.T, c *Client) {
				assert.Equal(t, int64(1024), c.maxBodySize)
			},
		},
		{
			name: "WithProgressFunc",
			opts: []Option{WithProgressFunc(func(done, total int64) {})},
			check: func(t *testing.T, c *Client) {
				assert.NotNil(t, c.progress)

				var called bool
				fn := func(done, total int64) {
					called = true
				}
				c.progress = fn
				c.progress(100, 200)
				assert.True(t, called)
			},
		},
		{
			name: "WithHTTPClient",
			opts: []Option{WithHTTPClient(customClient)},
			check: func(t *testing.T, c *Client) {
				assert.Equal(t, customClient, c.client)
				assert.Equal(t, 123*time.Second, c.client.Timeout)
			},
		},
		{
			name: "MultipleOptions",
			opts: []Option{
				WithTimeout(42 * time.Second),
				WithUserAgent("multi-agent"),
				WithMaxBodySize(2048),
			},
			check: func(t *testing.T, c *Client) {
				assert.Equal(t, 42*time.Second, c.timeout)
				assert.Equal(t, 42*time.Second, c.client.Timeout)
				assert.Equal(t, "multi-agent", c.userAgent)
				assert.Equal(t, int64(2048), c.maxBodySize)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New(tt.opts...)
			tt.check(t, c)
		})
	}
}
