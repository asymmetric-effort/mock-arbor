package arbor

import "testing"

func TestTerminalLength(t *testing.T) {
    out, _ := Dispatch("terminal length")
    if out == "" || out[:6] != "Usage:" {
        t.Fatalf("expected usage for missing arg; got %q", out)
    }
    out, _ = Dispatch("terminal length x")
    if out != "Invalid terminal length\n" {
        t.Fatalf("expected invalid length; got %q", out)
    }
    out, _ = Dispatch("terminal length 50")
    if out == "" || out[:8] != "Terminal" {
        t.Fatalf("expected set confirmation; got %q", out)
    }
}

