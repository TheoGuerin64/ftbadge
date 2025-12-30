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
	url := c.cdnBaseURL + endpoint

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create HTTP GET request for URL %q: %w", url, err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute HTTP GET request for URL %q: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received unexpected HTTP status code %d (%s) for URL %q", resp.StatusCode, resp.Status, url)
	}

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error decoding image from response body for URL %q: %w", url, err)
	}

	return img, nil
}

func (c *Client) Get(ctx context.Context, endpoint string, headers http.Header) (*http.Response, error) {
	url := c.apiBaseURL + endpoint

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create HTTP GET request for URL %q: %w", url, err)
	}

	for key, values := range headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute HTTP GET request for URL %q: %w", url, err)
	}

	return resp, nil
}

func (c *Client) PostForm(ctx context.Context, endpoint string, headers http.Header, data url.Values) (*http.Response, error) {
	url := c.apiBaseURL + endpoint
	body := strings.NewReader(data.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return nil, fmt.Errorf("unable to create HTTP POST request for URL %q: %w", url, err)
	}

	for key, values := range headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute HTTP POST request for URL %q: %w", url, err)
	}

	return resp, nil
}
