package upstream

import (
	"net"
	"net/http"
	"time"
)

// Default configuration values used when no option is provided.
const (
	defaultConnectionTimeout = 10 * time.Second
	defaultTimeout           = 30 * time.Second
	defaultUserAgent         = "tlgc/1.0"
	defaultMaxBodySize       = 1 << 30 // 1 GiB
)

// Client is an HTTP download client with configurable timeouts,
// size limits, progress reporting, and SHA256 verification.
type Client struct {
	// client is the underlying HTTP client used to make requests.
	client *http.Client
	// connectionTimeout is the timeout for dialing a TCP connection.
	connectionTimeout time.Duration
	// timeout is the overall timeout for an HTTP request.
	timeout time.Duration
	// userAgent is the value sent in the User-Agent header.
	userAgent string
	// maxBodySize is the maximum number of bytes read from the response body.
	maxBodySize int64
	// progress is an optional callback invoked periodically with download progress.
	progress func(done, total int64)
}

// Option configures a Client.
type Option func(*Client)

// WithHTTPClient sets the underlying HTTP client used for requests.
func WithHTTPClient(client *http.Client) Option {
	return func(c *Client) {
		c.client = client
	}
}

// WithConnectionTimeout sets the timeout for dialing a TCP connection.
func WithConnectionTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.connectionTimeout = timeout
	}
}

// WithTimeout sets the overall timeout for an HTTP request.
func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.timeout = timeout
	}
}

// WithMaxBodySize sets the maximum number of bytes read from the response body.
func WithMaxBodySize(n int64) Option {
	return func(c *Client) {
		c.maxBodySize = n
	}
}

// WithUserAgent sets the value sent in the User-Agent header.
func WithUserAgent(userAgent string) Option {
	return func(c *Client) {
		c.userAgent = userAgent
	}
}

// WithProgressFunc sets a callback that is invoked periodically with download progress.
func WithProgressFunc(fn func(done, total int64)) Option {
	return func(c *Client) {
		c.progress = fn
	}
}

// New creates a new Client with the given options.
func New(opts ...Option) *Client {
	c := &Client{
		connectionTimeout: defaultConnectionTimeout,
		timeout:           defaultTimeout,
		userAgent:         defaultUserAgent,
		maxBodySize:       defaultMaxBodySize,
	}

	for _, o := range opts {
		o(c)
	}

	if c.client == nil {
		c.client = &http.Client{
			Timeout: c.timeout,
			Transport: &http.Transport{
				DialContext: (&net.Dialer{
					Timeout: c.connectionTimeout,
				}).DialContext,
			},
		}
	}

	return c
}
