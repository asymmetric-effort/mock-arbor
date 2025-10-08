package arbor

import (
	"fmt"
	"mock-arbor/src/util"
	"strings"
)

// Dispatch takes a command line and returns (output, exitCode).
// It implements a small set of TMS-like commands with static responses for tooling.
func Dispatch(line string) (string, int) {
    cmd := util.NormalizeSpaces(line)
    switch {
	case util.EqualCmd(cmd, "help"), util.EqualCmd(cmd, "?"):
		return Help(), 0

	case strings.HasPrefix(cmd, "show "):
		return handleShow(cmd), 0

	case util.EqualCmd(cmd, "configure terminal"),
		util.EqualCmd(cmd, "conf t"):
		return "Entering configuration mode (emulated). Type 'exit' to leave.\n", 0

	case util.EqualCmd(cmd, "write memory"),
		util.EqualCmd(cmd, "wr mem"),
		util.EqualCmd(cmd, "copy running-config startup-config"):
		return "Configuration saved (emulated).\n", 0

    case strings.HasPrefix(cmd, "mitigation "):
        return handleMitigation(cmd), 0

    case strings.HasPrefix(cmd, "services "):
        return handleServices(cmd), 0

    case util.EqualCmd(cmd, "clear counters"):
        return "All counters cleared (emulated).\n", 0

	default:
		return fmt.Sprintf("%% Unrecognized command: %q\nType 'help' for a list of supported commands.\n", line), 1
	}
}
