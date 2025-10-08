package arbor

import (
    "strings"
    "testing"
)

func TestServicesSpBackupCreateFull(t *testing.T) {
    out, code := Dispatch("services sp backup create full")
    if code != 0 {
        t.Fatalf("expected code 0; got %d, out=%q", code, out)
    }
    if !strings.Contains(out, "SP backup: full") || !strings.Contains(out, "Status:   completed") {
        t.Fatalf("unexpected output: %q", out)
    }
}

func TestServicesUnknown(t *testing.T) {
    out, code := Dispatch("services sp unknown")
    if code != 0 || !strings.Contains(out, "Unrecognized 'services' command") {
        t.Fatalf("expected usage for unknown services; code=%d out=%q", code, out)
    }
}

