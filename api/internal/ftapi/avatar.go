package ftapi

import (
	"context"
	"fmt"

	"ftbadge/internal/cache"
	"ftbadge/internal/utils"
)

const (
	jpegQuality = 70
)

func (c *Client) GetAvatar(ctx context.Context, cm *cache.CacheManager, endpoint string) (string, error) {
	if cachedValue, isCached := cm.Get(cache.CacheKeyAvatar); isCached {
		return cachedValue, nil
	}

	image, err := c.fetchAndDecodeImage(ctx, endpoint)
	if err != nil {
		return "", fmt.Errorf("failed to fetch and decode avatar: %w", err)
	}
	image = utils.CropToSquare(image)

	jpegData, err := utils.EncodeToJPEG(image, jpegQuality)
	if err != nil {
		return "", fmt.Errorf("failed to encode avatar to JPEG from endpoint %q: %w", endpoint, err)
	}

	base64Image, err := utils.JPEGBytesToDataURI(jpegData)
	if err != nil {
		return "", fmt.Errorf("failed to convert JPEG bytes to base64 data URI from endpoint %q: %w", endpoint, err)
	}

	if err := cm.Set(cache.CacheKeyAvatar, base64Image); err != nil {
		return "", fmt.Errorf("failed to cache avatar: %w", err)
	}

	return base64Image, nil
}
