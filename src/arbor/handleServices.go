package arbor

import (
    "fmt"
    "mock-arbor/src/util"
)

// handleServices handles "services ..." commands used by Sightline TMS.
func handleServices(cmd string) string {
    switch {
    case util.EqualCmd(cmd, "services sp backup create full"):
        return getFixture("services/sp/backup_create_full", nil, "Sightline Services (emulated)\nSP backup: full\nStatus:   completed\nOutput:   /var/tmp/sightline/sp-backups/backup-full.tar.gz\n\n")
    default:
        return fmt.Sprintf("%% Unrecognized 'services' command: %q\nUsage (emulated):\n  services sp backup create full\n\n", cmd)
    }
}
