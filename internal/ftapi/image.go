package ftapi

import (
	"context"
	"encoding/base64"
	"fmt"
	"ftbadge/internal/cache"
	"io"
	"net/http"
	"time"
)

const (
	imageBase64CacheKeyPrefix = "image:base64:"
	imageBase64CacheTTL       = 24 * time.Hour * 7
)

func FetchImageAsBase64(ctx context.Context, imageURL string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, imageURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request for URL %q: %w", imageURL, err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch image from URL %q: %w", imageURL, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch image from URL %q: received unexpected status code %d (%s)", imageURL, resp.StatusCode, resp.Status)
	}

	imageBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read image data from response body for URL %q: %w", imageURL, err)
	}

	encodedImage := base64.StdEncoding.EncodeToString(imageBytes)
	return "data:image/png;base64," + encodedImage, nil
}

func GetOrCacheImageBase64(ctx context.Context, cs cache.CacheStore, imageURL string) (string, error) {
	cacheKey := imageBase64CacheKeyPrefix + base64.StdEncoding.EncodeToString([]byte(imageURL))

	cachedValue, found, err := cs.Get(ctx, cacheKey)
	if err != nil {
		return "", fmt.Errorf("error retrieving cached image for URL %q: %w", imageURL, err)
	}

	if !found {
		base64Image, err := FetchImageAsBase64(ctx, imageURL)
		if err != nil {
			return "", fmt.Errorf("could not encode image from URL %q: %w", imageURL, err)
		}

		if err = cs.Set(ctx, cacheKey, base64Image, imageBase64CacheTTL); err != nil {
			return "", fmt.Errorf("failed to cache base64 image for URL %q: %w", imageURL, err)
		}

		return base64Image, nil
	}

	return cachedValue, nil
}
