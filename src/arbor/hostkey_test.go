package arbor

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadOrGenerateHostKey(t *testing.T) {
	dir := t.TempDir()
	keyPath := filepath.Join(dir, "hostkey.pem")

	// First call should generate and persist.
	signer1, persisted1, err := LoadOrGenerateHostKey(keyPath)
	if err != nil {
		t.Fatalf("LoadOrGenerateHostKey (gen) error: %v", err)
	}
	if signer1 == nil || !persisted1 {
		t.Fatalf("expected signer and persisted=true on first generate")
	}
	if _, err := os.Stat(keyPath); err != nil {
		t.Fatalf("expected key file to exist: %v", err)
	}

	// Second call should load existing.
	signer2, persisted2, err := LoadOrGenerateHostKey(keyPath)
	if err != nil {
		t.Fatalf("LoadOrGenerateHostKey (load) error: %v", err)
	}
	if signer2 == nil || !persisted2 {
		t.Fatalf("expected persisted=true on load")
	}
	if string(signer1.PublicKey().Marshal()) != string(signer2.PublicKey().Marshal()) {
		t.Fatalf("public keys differ between generate and load")
	}
}

func TestLoadOrGenerateHostKeyEphemeralOnDir(t *testing.T) {
	dir := t.TempDir()
	// Create a directory where a file is expected; this forces persist to fail.
	keyPath := filepath.Join(dir, "as-dir")
	if err := os.Mkdir(keyPath, 0o700); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	signer, persisted, err := LoadOrGenerateHostKey(keyPath)
	if err != nil {
		t.Fatalf("LoadOrGenerateHostKey(dir) error: %v", err)
	}
	if signer == nil || persisted {
		t.Fatalf("expected signer and persisted=false when path is a directory")
	}
}
