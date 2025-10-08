package util

// IsAny returns true if s equals any of the candidates (case-insensitive, space-normalized).
func IsAny(s string, candidates ...string) bool {
	for _, c := range candidates {
		if EqualCmd(s, c) {
			return true
		}
	}
	return false
}
