package ftapi

import (
	"context"
	"fmt"
	"image"
	"net/http"
	"net/url"
	"strings"
	"time"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

type Client struct {
	client     *http.Client
	apiBaseURL string
	cdnBaseURL string
}

func NewClient(apiBaseURL string, cdnBaseURL string) *Client {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	return &Client{client, apiBaseURL, cdnBaseURL}
}

func (c *Client) fetchAndDecodeImage(ctx context.Context, endpoint string) (image.Image, error) {
	fullURL, err := url.JoinPath(c.cdnBaseURL, endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to construct URL from base %q and endpoint %q: %w", c.cdnBaseURL, endpoint, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create HTTP GET request for URL %q: %w", fullURL, err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute HTTP GET request for URL %q: %w", fullURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received unexpected HTTP status code %d (%s) for URL %q", resp.StatusCode, resp.Status, fullURL)
	}

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error decoding image from response body for URL %q: %w", fullURL, err)
	}

	return img, nil
}

func (c *Client) Get(ctx context.Context, endpoint string, headers http.Header) (*http.Response, error) {
	fullURL, err := url.JoinPath(c.apiBaseURL, endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to construct URL from base %q and endpoint %q: %w", c.apiBaseURL, endpoint, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create HTTP GET request for URL %q: %w", fullURL, err)
	}

	for key, values := range headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute HTTP GET request for URL %q: %w", fullURL, err)
	}

	return resp, nil
}

func (c *Client) PostForm(ctx context.Context, endpoint string, headers http.Header, data url.Values) (*http.Response, error) {
	fullURL, err := url.JoinPath(c.apiBaseURL, endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to construct URL from base %q and endpoint %q: %w", c.apiBaseURL, endpoint, err)
	}
	body := strings.NewReader(data.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fullURL, body)
	if err != nil {
		return nil, fmt.Errorf("unable to create HTTP POST request for URL %q: %w", fullURL, err)
	}

	for key, values := range headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute HTTP POST request for URL %q: %w", fullURL, err)
	}

	return resp, nil
}
