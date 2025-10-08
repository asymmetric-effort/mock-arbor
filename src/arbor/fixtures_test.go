package arbor

import (
    "os"
    "path/filepath"
    "strings"
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

func TestFixtureBadTemplateFallsBackToRaw(t *testing.T) {
    dir := t.TempDir()
    // Provide a bad template for system health
    path := filepath.Join(dir, "system", "health.txt")
    if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
        t.Fatalf("mkdir: %v", err)
    }
    bad := "Bad template {{ .UPTIME" // missing closing braces
    if err := os.WriteFile(path, []byte(bad), 0o644); err != nil {
        t.Fatalf("write fixture: %v", err)
    }
    old := os.Getenv("TMS_FIXTURES_DIR")
    t.Cleanup(func() { _ = os.Setenv("TMS_FIXTURES_DIR", old) })
    _ = os.Setenv("TMS_FIXTURES_DIR", dir)

    out, _ := Dispatch("show health")
    if strings.TrimSpace(out) != bad {
        t.Fatalf("expected raw bad template to be returned; got %q", out)
    }
}
