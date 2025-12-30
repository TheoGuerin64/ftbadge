package utils

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
)

func CropToSquare(img image.Image) image.Image {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	size := min(width, height)

	xOffset := (width - size) / 2
	yOffset := (height - size) / 2
	cropRect := image.Rect(xOffset, yOffset, xOffset+size, yOffset+size)

	square := image.NewRGBA(image.Rect(0, 0, size, size))
	draw.Draw(square, square.Bounds(), img, cropRect.Min, draw.Src)

	return square
}

func EncodeToJPEG(img image.Image, quality int) ([]byte, error) {
	opt := jpeg.Options{
		Quality: quality,
	}

	buf := new(bytes.Buffer)
	if err := jpeg.Encode(buf, img, &opt); err != nil {
		return nil, fmt.Errorf("JPEG encoding failed: %w", err)
	}
	return buf.Bytes(), nil
}

func JPEGBytesToDataURI(jpegData []byte) (string, error) {
	base64Data := base64.StdEncoding.EncodeToString(jpegData)
	return "data:image/jpeg;base64," + base64Data, nil
}
