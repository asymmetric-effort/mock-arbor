package arbor

import (
	"bufio"
	"net"
	"strings"
	"testing"
	"time"

	gliderssh "github.com/gliderlabs/ssh"
	xssh "golang.org/x/crypto/ssh"
)

// startTestServer starts an SSH server bound to a random localhost port and
// returns the address and a cleanup function.
func startTestServer(t *testing.T) (addr string, cleanup func()) {
	t.Helper()

	signer, _, err := LoadOrGenerateHostKey(t.TempDir() + "/hostkey.pem")
	if err != nil {
		t.Fatalf("hostkey: %v", err)
	}

	srv := &gliderssh.Server{
		PasswordHandler: func(ctx gliderssh.Context, pass string) bool {
			return ctx.User() == "tms" && pass == "tms123!"
		},
		Banner:        Banner,
		BannerHandler: func(ctx gliderssh.Context) string { return Banner },
		Handler:       HandleSession,
	}
	srv.AddHostKey(signer)

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	done := make(chan struct{})
	go func() {
		_ = srv.Serve(ln)
		close(done)
	}()

	cleanup = func() {
		_ = ln.Close()
		select {
		case <-done:
		case <-time.After(2 * time.Second):
			// give the server a moment to stop
		}
	}
	return ln.Addr().String(), cleanup
}

func dialClient(t *testing.T, addr string) *xssh.Client {
	t.Helper()
	cfg := &xssh.ClientConfig{
		User:            "tms",
		Auth:            []xssh.AuthMethod{xssh.Password("tms123!")},
		HostKeyCallback: xssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}
	c, err := xssh.Dial("tcp", addr, cfg)
	if err != nil {
		t.Fatalf("ssh dial: %v", err)
	}
	return c
}

func readUntilFromChan(t *testing.T, dataCh <-chan []byte, errCh <-chan error, substr string, timeout time.Duration) string {
	t.Helper()
	deadline := time.Now().Add(timeout)
	var b strings.Builder
	for time.Now().Before(deadline) {
		select {
		case p := <-dataCh:
			b.Write(p)
			if strings.Contains(b.String(), substr) {
				return b.String()
			}
		case <-time.After(50 * time.Millisecond):
		case <-errCh:
			t.Fatalf("stream closed waiting for %q; got %q", substr, b.String())
		}
	}
	t.Fatalf("timeout waiting for %q; got %q", substr, b.String())
	return ""
}

func TestSSHExecCommands(t *testing.T) {
	addr, cleanup := startTestServer(t)
	defer cleanup()

	c := dialClient(t, addr)
	defer c.Close()

	// show version via exec
	sess, err := c.NewSession()
	if err != nil {
		t.Fatalf("new session: %v", err)
	}
	out, err := sess.CombinedOutput("show version")
	_ = sess.Close()
	if err != nil {
		t.Fatalf("exec show version error: %v; out=%q", err, string(out))
	}
	if !strings.Contains(string(out), "TMS Software (emulated)") {
		t.Fatalf("unexpected output: %q", out)
	}

	// mitigation start via exec
	sess2, err := c.NewSession()
	if err != nil {
		t.Fatalf("new session2: %v", err)
	}
	out, err = sess2.CombinedOutput("mitigation start")
	_ = sess2.Close()
    if err != nil || !strings.Contains(string(out), "Mitigation started (emulated)") {
        t.Fatalf("mitigation start failed: err=%v out=%q", err, string(out))
    }

    // unknown command should return non-zero
    sess3, err := c.NewSession()
    if err != nil { t.Fatalf("new session3: %v", err) }
    out, err = sess3.CombinedOutput("bogus-cmd-xyz")
    _ = sess3.Close()
    if err == nil || !strings.Contains(string(out), "Unrecognized command") {
        t.Fatalf("expected exec error and message; err=%v out=%q", err, string(out))
    }
}

func TestSSHInteractiveHelpAndExit(t *testing.T) {
	addr, cleanup := startTestServer(t)
	defer cleanup()

	c := dialClient(t, addr)
	defer c.Close()

	sess, err := c.NewSession()
	if err != nil {
		t.Fatalf("new session: %v", err)
	}
	defer sess.Close()

	// Request PTY and start shell
	if err := sess.RequestPty("xterm", 80, 24, xssh.TerminalModes{}); err != nil {
		t.Fatalf("pty: %v", err)
	}

	stdout, err := sess.StdoutPipe()
	if err != nil {
		t.Fatalf("stdout: %v", err)
	}
	stdin, err := sess.StdinPipe()
	if err != nil {
		t.Fatalf("stdin: %v", err)
	}

	if err := sess.Shell(); err != nil {
		t.Fatalf("shell: %v", err)
	}

	r := bufio.NewReader(stdout)
	dataCh := make(chan []byte, 16)
	errCh := make(chan error, 1)
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := r.Read(buf)
			if n > 0 {
				cp := make([]byte, n)
				copy(cp, buf[:n])
				dataCh <- cp
			}
			if err != nil {
				errCh <- err
				return
			}
		}
	}()

	// Expect banner then prompt
	out := readUntilFromChan(t, dataCh, errCh, Prompt, 3*time.Second)
	norm := strings.ReplaceAll(out, "\r\n", "\n")
	if !strings.Contains(norm, Banner) {
		t.Fatalf("missing banner; got %q", out)
	}

    // Send blank line; expect prompt again
    if _, err := stdin.Write([]byte("\n")); err != nil {
        t.Fatalf("write blank: %v", err)
    }
    out = readUntilFromChan(t, dataCh, errCh, Prompt, 2*time.Second)

    // Send help and expect prompt again
    if _, err := stdin.Write([]byte("help\n")); err != nil {
        t.Fatalf("write help: %v", err)
    }
    out = readUntilFromChan(t, dataCh, errCh, Prompt, 3*time.Second)
    norm = strings.ReplaceAll(out, "\r\n", "\n")
    if !strings.Contains(norm, "Available commands (emulated)") {
        t.Fatalf("missing help output; got %q", out)
    }

	// Exit
	if _, err := stdin.Write([]byte("exit\n")); err != nil {
		t.Fatalf("write exit: %v", err)
	}
	// Expect bye and session to close
	_ = readUntilFromChan(t, dataCh, errCh, "Bye.", 2*time.Second)
	if err := sess.Wait(); err != nil {
		// Close returns io.EOF sometimes; treat as success if not clearly failed
	}
}
