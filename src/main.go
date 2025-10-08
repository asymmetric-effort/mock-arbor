// Command tms_ssh_emulator runs a minimal SSH server that emulates an
// Arbor TMS-style CLI (prompt, a handful of "show" commands, and exec support).
//
// Usage:
//
//	go mod init example.com/tms_ssh_emulator
//	go get github.com/gliderlabs/ssh@v0.3.6 golang.org/x/crypto/ssh
//	go build -o tms_ssh_emulator
//	./tms_ssh_emulator
//
// Env vars:
//
//	TMS_SSH_ADDR         (default ":2222")
//	TMS_USERNAME         (default "tms")
//	TMS_PASSWORD         (default "tms123!")
//	TMS_HOSTKEY_FILE     (default "hostkey.pem")
package main

import (
	"errors"
	gliderssh "github.com/gliderlabs/ssh"
	"log"
	"mock-arbor/src/arbor"
	"mock-arbor/src/util"
	"net"
)

const (
	defaultAddr       = ":2222"
	defaultUser       = "tms"
	defaultPass       = "tms123!"
	defaultHostKeyPem = "hostkey.pem"
)

func main() {
	addr := util.GetEnv("TMS_SSH_ADDR", defaultAddr)
	user := util.GetEnv("TMS_USERNAME", defaultUser)
	pass := util.GetEnv("TMS_PASSWORD", defaultPass)
	hostKeyPath := util.GetEnv("TMS_HOSTKEY_FILE", defaultHostKeyPem)

	signer, persisted, err := arbor.LoadOrGenerateHostKey(hostKeyPath)
	if err != nil {
		log.Fatalf("host key error: %v", err)
	}
	if persisted {
		log.Printf("loaded host key: %s", hostKeyPath)
	} else {
		log.Printf("using ephemeral host key (no %s found)", hostKeyPath)
	}

	server := gliderssh.Server{
		Addr: addr,
		// Password auth only (for demo). Consider adding pubkey auth if needed.
		PasswordHandler: func(ctx gliderssh.Context, candidatePass string) bool {
			return ctx.User() == user && candidatePass == pass
		},
		Banner: arbor.Banner,
		BannerHandler: func(ctx gliderssh.Context) string {
			return arbor.Banner
		},
		Handler: func(s gliderssh.Session) {
			arbor.HandleSession(s)
		},
	}
	server.AddHostKey(signer)

	log.Printf("TMS SSH emulator listening on %s (user=%q)", addr, user)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, net.ErrClosed) {
		log.Fatalf("server error: %v", err)
	}
}
