package utils

import (
	"compress/gzip"
	"fmt"
	"io"
)

func DecompressGzip(input io.Reader) ([]byte, error) {
	gzipReader, err := gzip.NewReader(input)
	if err != nil {
		return nil, fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzipReader.Close()

	data, err := io.ReadAll(gzipReader)
	if err != nil {
		return nil, fmt.Errorf("error reading decompressed gzip data: %w", err)
	}

	return data, nil
}
