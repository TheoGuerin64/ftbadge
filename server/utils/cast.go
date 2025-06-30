package utils

import "fmt"

func AnyToStringPointer(value any) (*string, error) {
	if value == nil {
		return nil, nil
	}
	str, ok := value.(string)
	if !ok {
		return nil, fmt.Errorf("expected string but got %T", value)
	}
	return &str, nil
}
