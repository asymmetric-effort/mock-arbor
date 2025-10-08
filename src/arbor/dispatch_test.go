package arbor

import (
	"strings"
	"testing"
)

func TestHelpAndUnknown(t *testing.T) {
	out, code := Dispatch("help")
	if code != 0 || !strings.Contains(out, "Available commands (emulated)") {
		t.Fatalf("help dispatch failed: code=%d out=%q", code, out)
	}

	out, code = Dispatch("?")
	if code != 0 || !strings.Contains(out, "Available commands (emulated)") {
		t.Fatalf("? dispatch failed: code=%d out=%q", code, out)
	}

	out, code = Dispatch("not-a-command")
	if code == 0 || !strings.Contains(out, "Unrecognized command") {
		t.Fatalf("unknown should be error: code=%d out=%q", code, out)
	}
}

func TestShowAndMitigation(t *testing.T) {
	out, code := Dispatch("show version")
	if code != 0 || !strings.Contains(out, "TMS Software (emulated)") {
		t.Fatalf("show version failed: code=%d out=%q", code, out)
	}

	out, code = Dispatch("show interfaces")
	if code != 0 || !strings.Contains(out, "Interface Summary (emulated)") {
		t.Fatalf("show interfaces failed: %q", out)
	}

	out, code = Dispatch("mitigation start")
	if code != 0 || !strings.Contains(out, "Mitigation started (emulated)") {
		t.Fatalf("mitigation start failed: %q", out)
	}

	out, code = Dispatch("mitigation apply pol1")
	if code != 0 || !strings.Contains(out, "Applied mitigation policy \"pol1\"") {
		t.Fatalf("mitigation apply failed: %q", out)
	}
}
