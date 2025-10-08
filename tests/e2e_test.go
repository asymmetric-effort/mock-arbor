package tests

import (
    "bufio"
    "context"
    "net"
    "os"
    "os/exec"
    "path/filepath"
    "runtime"
    "strings"
    "testing"
    "time"

    xssh "golang.org/x/crypto/ssh"
)

func findFreePort(t *testing.T) string {
    t.Helper()
    ln, err := net.Listen("tcp", "127.0.0.1:0")
    if err != nil {
        t.Fatalf("listen 0: %v", err)
    }
    addr := ln.Addr().String()
    _ = ln.Close()
    return addr
}

func waitForTCP(t *testing.T, addr string, timeout time.Duration) {
    t.Helper()
    deadline := time.Now().Add(timeout)
    for time.Now().Before(deadline) {
        c, err := net.DialTimeout("tcp", addr, 200*time.Millisecond)
        if err == nil {
            _ = c.Close()
            return
        }
        time.Sleep(100 * time.Millisecond)
    }
    t.Fatalf("server did not start on %s", addr)
}

func buildBinary(t *testing.T, outPath string) {
    t.Helper()
    root := findModuleRoot(t)
    cmd := exec.Command("go", "build", "-o", outPath, "./src")
    cmd.Dir = root
    cmd.Env = append(os.Environ(), "CGO_ENABLED=0")
    out, err := cmd.CombinedOutput()
    if err != nil {
        t.Fatalf("go build failed: %v\n%s", err, string(out))
    }
}

func findModuleRoot(t *testing.T) string {
    t.Helper()
    dir, err := os.Getwd()
    if err != nil { t.Fatalf("getwd: %v", err) }
    for i := 0; i < 5; i++ {
        if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
            return dir
        }
        parent := filepath.Dir(dir)
        if parent == dir {
            break
        }
        dir = parent
    }
    t.Fatalf("could not locate go.mod above %s", dir)
    return ""
}

func startServer(t *testing.T, addr string) (*exec.Cmd, func()) {
    t.Helper()
    tmp := t.TempDir()
    bin := filepath.Join(tmp, "tms_ssh_emulator")
    if runtime.GOOS == "windows" { bin += ".exe" }
    buildBinary(t, bin)

    hostkey := filepath.Join(tmp, "hostkey.pem")
    cmd := exec.Command(bin)
    cmd.Env = append(os.Environ(),
        "TMS_SSH_ADDR="+addr,
        "TMS_USERNAME=tms",
        "TMS_PASSWORD=tms123!",
        "TMS_HOSTKEY_FILE="+hostkey,
    )
    // Capture logs for debugging
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    if err := cmd.Start(); err != nil {
        t.Fatalf("start server: %v", err)
    }

    cleanup := func() {
        ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
        defer cancel()
        _ = cmd.Process.Signal(os.Interrupt)
        done := make(chan struct{})
        go func() { _ = cmd.Wait(); close(done) }()
        select {
        case <-done:
        case <-ctx.Done():
            _ = cmd.Process.Kill()
        }
    }
    return cmd, cleanup
}

func sshClient(t *testing.T, addr string) *xssh.Client {
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

func TestE2E_ExecAndInteractive(t *testing.T) {
    addr := findFreePort(t)
    cmd, cleanup := startServer(t, addr)
    defer cleanup()
    _ = cmd
    waitForTCP(t, addr, 5*time.Second)

    c := sshClient(t, addr)
    defer c.Close()

    // Exec: show version
    sess, err := c.NewSession()
    if err != nil { t.Fatalf("new session: %v", err) }
    out, err := sess.CombinedOutput("show version")
    _ = sess.Close()
    if err != nil || !strings.Contains(string(out), "TMS Software") {
        t.Fatalf("exec show version failed: err=%v out=%q", err, string(out))
    }

    // Exec: services sp backup create full
    sess2, err := c.NewSession()
    if err != nil { t.Fatalf("new session2: %v", err) }
    out, err = sess2.CombinedOutput("services sp backup create full")
    _ = sess2.Close()
    if err != nil || !strings.Contains(string(out), "backup") {
        t.Fatalf("exec services backup failed: err=%v out=%q", err, string(out))
    }

    // Interactive session: help then exit
    sess3, err := c.NewSession()
    if err != nil { t.Fatalf("new session3: %v", err) }
    defer sess3.Close()
    if err := sess3.RequestPty("xterm", 80, 24, xssh.TerminalModes{}); err != nil {
        t.Fatalf("pty: %v", err)
    }
    stdout, err := sess3.StdoutPipe()
    if err != nil { t.Fatalf("stdout: %v", err) }
    stdin, err := sess3.StdinPipe()
    if err != nil { t.Fatalf("stdin: %v", err) }
    if err := sess3.Shell(); err != nil { t.Fatalf("shell: %v", err) }

    r := bufio.NewReader(stdout)
    // Read banner+prompt
    buf := make([]byte, 4096)
    n, _ := r.Read(buf)
    if !strings.Contains(string(buf[:n]), "tms#") {
        t.Fatalf("missing prompt in interactive start: %q", string(buf[:n]))
    }
    // Send help, expect menu
    if _, err := stdin.Write([]byte("help\n")); err != nil { t.Fatalf("write help: %v", err) }
    time.Sleep(100 * time.Millisecond)
    n, _ = r.Read(buf)
    if !strings.Contains(string(buf[:n]), "Available commands") {
        t.Fatalf("missing help output: %q", string(buf[:n]))
    }
    // Exit
    if _, err := stdin.Write([]byte("exit\n")); err != nil { t.Fatalf("write exit: %v", err) }
}
