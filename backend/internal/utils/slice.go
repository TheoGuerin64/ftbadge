package utils

import (
	"fmt"
)

func MapSlice[Input any, Output any](elements []Input, transform func(Input) (Output, error)) ([]Output, error) {
	mapped := make([]Output, len(elements))
	for index, element := range elements {
		transformed, err := transform(element)
		if err != nil {
			return nil, fmt.Errorf("failed to transform element at index %d: %w", index, err)
		}
		mapped[index] = transformed
	}
	return mapped, nil
}
