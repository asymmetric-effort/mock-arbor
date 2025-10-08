package arbor

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"golang.org/x/crypto/ssh"
	"os"
	"path/filepath"
)

// LoadOrGenerateHostKey - attempts to read an existing PEM-encoded private key
// from path. If missing, it generates an RSA key and returns an ssh.Signer.
// It returns (signer, persisted, error), where persisted indicates whether a
// file was successfully loaded (true) or a key was generated in-memory (false).
func LoadOrGenerateHostKey(path string) (ssh.Signer, bool, error) {
	if data, err := os.ReadFile(path); err == nil {
		signer, err := ssh.ParsePrivateKey(data)
		if err != nil {
			return nil, false, fmt.Errorf("parse hostkey %s: %w", path, err)
		}
		return signer, true, nil
	}

	// Generate a new RSA key (ephemeral unless we persist it).
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, false, fmt.Errorf("generate rsa: %w", err)
	}
	signer, err := ssh.NewSignerFromKey(key)
	if err != nil {
		return nil, false, fmt.Errorf("signer: %w", err)
	}

	// Best-effort: persist if directory is writable and file doesn't exist.
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o700); err == nil {
		f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
		if err == nil {
			defer func(f *os.File) {
				if cErr := f.Close(); cErr != nil {
					panic(cErr)
				}
			}(f)
			privBytes := x509.MarshalPKCS1PrivateKey(key)
			pemBlock := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: privBytes}
			if err := pem.Encode(f, pemBlock); err == nil {
				return signer, true, nil
			}
			// If encode fails, fall through to ephemeral.
		}
	}
	return signer, false, nil
}
