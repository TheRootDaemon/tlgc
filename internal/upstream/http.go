package upstream

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

func (c *Client) newRequest(ctx context.Context, url string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		url,
		nil,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", c.userAgent)

	return req, nil
}

func (c *Client) send(req *http.Request) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) validateResponse(resp *http.Response) error {
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		_ = resp.Body.Close()
		return fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	return nil
}

func (c *Client) wrapBody(
	body io.ReadCloser,
	contentLength int64,
) io.ReadCloser {
	if c.maxBodySize > 0 {
		body = &limitedBody{
			Reader: io.LimitReader(body, c.maxBodySize+1),
			closer: body,
			limit:  c.maxBodySize,
		}
	}

	if c.progress != nil {
		body = &progressReader{
			reader: body,
			total:  contentLength,
			fn:     c.progress,
		}
	}

	return body
}

func (c *Client) execute(ctx context.Context, url string) (*http.Response, error) {
	req, err := c.newRequest(ctx, url)
	if err != nil {
		return nil, err
	}

	resp, err := c.send(req)
	if err != nil {
		return nil, err
	}

	if err = c.validateResponse(resp); err != nil {
		return nil, err
	}

	resp.Body = c.wrapBody(resp.Body, resp.ContentLength)

	return resp, nil
}
