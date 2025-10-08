package arbor

import (
    "os"
    "path/filepath"
    "testing"
)

func TestFixtureOverride(t *testing.T) {
    dir := t.TempDir()
    // Provide a custom template for mitigation summary
    path := filepath.Join(dir, "mitigation", "summary.txt")
    if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
        t.Fatalf("mkdir: %v", err)
    }
    if err := os.WriteFile(path, []byte("Custom Summary Fixture\n"), 0o644); err != nil {
        t.Fatalf("write fixture: %v", err)
    }
    old := os.Getenv("TMS_FIXTURES_DIR")
    t.Cleanup(func() { _ = os.Setenv("TMS_FIXTURES_DIR", old) })
    _ = os.Setenv("TMS_FIXTURES_DIR", dir)

    out, _ := Dispatch("show mitigation summary")
    if out != "Custom Summary Fixture\n" {
        t.Fatalf("fixture not applied; got %q", out)
    }
}

