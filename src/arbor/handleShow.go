package arbor

import (
	"fmt"
	"mock-arbor/src/util"
	"strings"
	"time"
)

// handleShow returns output for "show ..." subcommands.
func handleShow(cmd string) string {
	switch {
	case util.EqualCmd(cmd, "show version"):
		return "TMS Software (emulated)\nVersion: 99.9.9-emu\nBuild: 2025-10-08\nUptime: " +
			time.Since(StartTime()).Round(time.Second).String() + "\n\n"

	case util.EqualCmd(cmd, "show running-config"), util.EqualCmd(cmd, "show run"):
		return `!
! Emulated running configuration
hostname tms
username admin role admin
logging host 192.0.2.10
ntp server 198.51.100.20
interface mgt0 ip address 10.0.0.10/24
!
mitigation default-action permit
policy default threshold 50000pps
!
end

`

	case util.EqualCmd(cmd, "show mitigation"), util.EqualCmd(cmd, "show mitigation status"):
		return `Mitigation Status (emulated)
--------------------------------
State:            IDLE
Active Policies:  0
Blocked Hosts:    0
Last Event:       none

`

	case strings.HasPrefix(cmd, "show policy "):
		pol := strings.TrimPrefix(cmd, "show policy ")
		return fmt.Sprintf("Policy %q (emulated)\n  Match: any\n  Action: rate-limit 1G\n  Threshold: 50000pps\n\n", pol)

	case util.EqualCmd(cmd, "show interfaces"), util.EqualCmd(cmd, "show int"):
		return `Interface Summary (emulated)
Name     State  Speed   RX pps   TX pps   RX bps     TX bps
------   -----  ------  -------  -------  ---------  ---------
mgt0     up     1G      120      95       128000     102400
mit0     up     10G     5200     4800     8.2e9      7.5e9

`

	case util.EqualCmd(cmd, "show logging"):
		return `Logging (emulated)
  Buffer size:    4096
  Remote server:  192.0.2.10:514
  Level:          info

`

	default:
		return fmt.Sprintf("%% Unrecognized 'show' command: %q\nTry: show version | show running-config | show mitigation | show interfaces\n\n", cmd)
	}
}
