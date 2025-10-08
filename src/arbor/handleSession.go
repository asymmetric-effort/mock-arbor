package arbor

import (
	"bufio"
	gliderssh "github.com/gliderlabs/ssh"
	"io"
	"mock-arbor/src/util"
	"strings"
)

// HandleSession routes between "exec" mode and interactive shell mode.
// It supports a tiny Arbor TMS-like command set sufficient for tooling tests.
func HandleSession(s gliderssh.Session) {
	// Best-effort close; ignore close errors (can be EOF after exit).
	defer func(s gliderssh.Session) { _ = s.Close() }(s)

	// Handle "exec" (single command) if present.
	if cmd := s.Command(); len(cmd) > 0 {
		out, code := Dispatch(strings.Join(cmd, " "))
		_, err := io.WriteString(s, out)
		if err != nil {
			return
		}
		if err = s.Exit(code); err != nil {
			return
		}
		return
	}

	// Interactive shell.
	width, height, ok := s.Pty()
	if ok {
		_ = width
		_ = height
	} // we don't need actual sizing for this minimal emulator

	if _, err := io.WriteString(s, Banner); err != nil {
		return
	}

	if _, err := io.WriteString(s, Prompt); err != nil {
		return
	}

	reader := bufio.NewReader(s)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			// EOF or client disconnected
			return
		}
		line = strings.TrimSpace(line)
		if line == "" {
			if _, err = io.WriteString(s, Prompt); err != nil {
				return
			}
			continue
		}
		// Basic exit commands.
		if util.IsAny(line, "exit", "quit", "logout") {
			if _, err = io.WriteString(s, "Bye.\n"); err != nil {
				return
			}
			if err = s.Exit(0); err != nil {
				return
			}
			return
		}
		// Dispatch and render.
		out, _ := Dispatch(line)
		if _, err = io.WriteString(s, out); err != nil {
			return
		}
		if _, err = io.WriteString(s, Prompt); err != nil {
			return
		}
	}
}
