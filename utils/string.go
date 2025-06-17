package utils

import (
	"fmt"
	"strconv"
)

func ParseStringToInt(s string) (int, error) {
	if s == "" {
		return 0, nil
	}
	value, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("failed to parse string to int: %w", err)
	}
	return value, nil
}

func ParseBoolToString(b *bool) string {
	if b == nil {
		return ""
	}
	if *b {
		return "true"
	}
	return "false"
}
