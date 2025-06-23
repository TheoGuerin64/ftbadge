package ftapi

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"ftbadge/internal/cache"
	"image"
	"image/draw"
	"image/jpeg"
	"net/http"
	"time"

	_ "image/gif"
	_ "image/png"
)

const (
	imageBase64CacheKeyPrefix = "image:base64:"
	imageBase64CacheTTL       = 24 * time.Hour * 7
	imageBase64JPEGQuality    = 70
)

func cropImageToSquare(img image.Image) image.Image {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	size := min(width, height)

	xOffset := (width - size) / 2
	yOffset := (height - size) / 2
	cropRect := image.Rect(xOffset, yOffset, xOffset+size, yOffset+size)

	squareImg := image.NewRGBA(image.Rect(0, 0, size, size))
	draw.Draw(squareImg, squareImg.Bounds(), img, cropRect.Min, draw.Src)

	return squareImg
}

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

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return "", fmt.Errorf("image decoding failed for URL %q: %w", imageURL, err)
	}

	squareImg := cropImageToSquare(img)

	buf := bytes.NewBuffer(nil)
	if err := jpeg.Encode(buf, squareImg, &jpeg.Options{Quality: imageBase64JPEGQuality}); err != nil {
		return "", fmt.Errorf("JPEG encoding failed for image from URL %q: %w", imageURL, err)
	}

	base64Data := base64.StdEncoding.EncodeToString(buf.Bytes())
	return "data:image/jpg;base64," + base64Data, nil
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
