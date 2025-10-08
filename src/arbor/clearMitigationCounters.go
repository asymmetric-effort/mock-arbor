package arbor

import (
	"fmt"
	"strings"
)

func clearMitigationCounters(cmd string) string {
	// clear mitigation counters [<id>]
	fields := strings.Fields(cmd)
	if len(fields) == 3 {
		return getFixture("mitigation/clear_all", nil, "Mitigation counters cleared (all) (emulated).\n")
	}
	id := fields[3]
	return getFixture("mitigation/clear_id", map[string]string{"ID": id}, fmt.Sprintf("Mitigation counters cleared for %s (emulated).\n", id))
}
