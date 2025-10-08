package arbor

import (
	"fmt"
	"mock-arbor/src/util"
	"strings"
)

// handleMitigation handles "mitigation ..." subcommands (emulated actions).
func handleMitigation(cmd string) string {
    switch {
    case util.EqualCmd(cmd, "mitigation start"):
        return getFixture("mitigation/start", nil, "Mitigation started (emulated). No actual traffic is affected.\n")
    case util.EqualCmd(cmd, "mitigation stop"):
        return getFixture("mitigation/stop", nil, "Mitigation stopped (emulated).\n")
    case strings.HasPrefix(cmd, "mitigation apply "):
        pol := strings.TrimPrefix(cmd, "mitigation apply ")
        return getFixture("mitigation/apply", map[string]string{"POLICY": pol}, fmt.Sprintf("Applied mitigation policy %q (emulated).\n", pol))
    case util.EqualCmd(cmd, "mitigation list"):
        return mitigationList()
    default:
        return "Usage (emulated):\n  mitigation start\n  mitigation stop\n  mitigation apply <policy>\n"
    }
}
