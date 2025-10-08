package util

// EqualCmd compares commands case-insensitively after normalizing spaces.
func EqualCmd(a, b string) bool {
	return NormalizeSpaces(a) == NormalizeSpaces(b)
}
