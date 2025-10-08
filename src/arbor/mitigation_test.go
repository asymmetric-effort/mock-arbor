package arbor

import (
	"strings"
	"testing"
)

func TestMitigationStopAndUsage(t *testing.T) {
	out, code := Dispatch("mitigation stop")
	if code != 0 || !strings.Contains(out, "Mitigation stopped") {
		t.Fatalf("mitigation stop failed: %q (code=%d)", out, code)
	}
	out, code = Dispatch("mitigation ?")
	if code != 0 || !strings.HasPrefix(out, "Usage") {
		t.Fatalf("mitigation usage failed: %q (code=%d)", out, code)
	}
}
