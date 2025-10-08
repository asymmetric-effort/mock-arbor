package util

import "strings"

// NormalizeSpaces trims and collapses runs of whitespace into single spaces.
func NormalizeSpaces(s string) string {
	fields := strings.Fields(strings.TrimSpace(s))
	return strings.ToLower(strings.Join(fields, " "))
}
