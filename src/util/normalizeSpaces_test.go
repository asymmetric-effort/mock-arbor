package util

import "testing"

func TestNormalizeSpaces(t *testing.T) {
	cases := map[string]string{
		"":                         "",
		"  ":                       "",
		"SHOW   VERSION":           "show version",
		"  mitigation   apply  p1": "mitigation apply p1",
		"\tConf\t t\t":             "conf t",
	}
	for in, want := range cases {
		got := NormalizeSpaces(in)
		if got != want {
			t.Fatalf("NormalizeSpaces(%q)=%q; want %q", in, got, want)
		}
	}
}

func TestEqualCmdAndIsAny(t *testing.T) {
	if !EqualCmd("  SHOW  VERSION ", "show version") {
		t.Fatalf("EqualCmd should match case/space-insensitively")
	}
	if IsAny("foo", "bar", "baz") {
		t.Fatalf("IsAny false positive")
	}
	if !IsAny("  WR   MEM  ", "write memory", "wr mem") {
		t.Fatalf("IsAny should match one of candidates")
	}
}
