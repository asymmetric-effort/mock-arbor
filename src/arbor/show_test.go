package arbor

import (
	"strings"
	"testing"
)

func TestShowRunningConfig(t *testing.T) {
	out, code := Dispatch("show running-config")
	if code != 0 || !strings.Contains(out, "Emulated running configuration") {
		t.Fatalf("show running-config failed: code=%d out=%q", code, out)
	}
}

func TestShowMitigationVariants(t *testing.T) {
	out, code := Dispatch("show mitigation")
	if code != 0 || !strings.Contains(out, "Mitigation Status (emulated)") {
		t.Fatalf("show mitigation failed: %q", out)
	}
	out, code = Dispatch("show mitigation status")
	if code != 0 || !strings.Contains(out, "Mitigation Status (emulated)") {
		t.Fatalf("show mitigation status failed: %q", out)
	}
}

func TestShowPolicyAndLogging(t *testing.T) {
	out, code := Dispatch("show policy alpha")
	if code != 0 || !strings.Contains(out, "Policy \"alpha\"") {
		t.Fatalf("show policy failed: %q", out)
	}
	out, code = Dispatch("show logging")
	if code != 0 || !strings.Contains(out, "Logging (emulated)") {
		t.Fatalf("show logging failed: %q", out)
	}
}

func TestShowUnknown(t *testing.T) {
	out, code := Dispatch("show nonsense")
	if code != 0 || !strings.Contains(out, "Unrecognized 'show' command") {
		t.Fatalf("show unknown failed: code=%d out=%q", code, out)
	}
}
