package arbor

import (
	"strings"
	"testing"
)

func TestConfigureAndWriteAndClear(t *testing.T) {
	out, code := Dispatch("configure terminal")
	if code != 0 || !strings.Contains(out, "Entering configuration mode") {
		t.Fatalf("configure terminal failed: %q", out)
	}
	out, code = Dispatch("wr mem")
	if code != 0 || !strings.Contains(out, "Configuration saved") {
		t.Fatalf("wr mem failed: %q", out)
	}
	out, code = Dispatch("clear counters")
	if code != 0 || !strings.Contains(out, "All counters cleared") {
		t.Fatalf("clear counters failed: %q", out)
	}
}
