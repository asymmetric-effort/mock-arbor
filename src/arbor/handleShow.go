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

    case util.EqualCmd(cmd, "show interfaces brief"):
        return `Interfaces (brief)
Name   State  Speed
mgt0   up     1G
mit0   up     10G

`

    case util.EqualCmd(cmd, "show logging"):
        return `Logging (emulated)
  Buffer size:    4096
  Remote server:  192.0.2.10:514
  Level:          info

`

    // Mitigation ops
    case util.EqualCmd(cmd, "show mitigation summary"):
        return getFixture("mitigation/summary", nil, `Mitigation Summary (emulated)
ID    Name        State   Packets Dropped  Bps Blocked
101   ddos-icmp   IDLE    0                0
102   ddos-syn    ACTIVE  1200000          350M

`)
    case strings.HasPrefix(cmd, "show mitigation detail "):
        id := strings.TrimPrefix(cmd, "show mitigation detail ")
        return getFixture("mitigation/detail", map[string]string{"ID": id}, fmt.Sprintf(`Mitigation Detail (emulated)
ID:            %s
Name:          example-%s
State:         ACTIVE
Targets:       203.0.113.10/32
Duration:      00:10:42
Packets Dropped: 1.2M
Bps Blocked:     350M

`, id, id))
    case util.EqualCmd(cmd, "show mitigation statistics"), util.EqualCmd(cmd, "show mitigation counters"):
        return getFixture("mitigation/counters", nil, `Mitigation Counters (emulated)
Total Active:  1
Total Idle:    1
Drops pps:     75000
Drops bps:     350M

`)
    case strings.HasPrefix(cmd, "show mitigation statistics ") || strings.HasPrefix(cmd, "show mitigation counters "):
        // accept "show mitigation statistics <id>" or counters form
        fields := strings.Fields(cmd)
        id := fields[len(fields)-1]
        return getFixture("mitigation/counters_id", map[string]string{"ID": id}, fmt.Sprintf(`Mitigation Counters (emulated)
ID: %s
Drops pps:  72000
Drops bps:  340M

`, id))

    // TMS groups / diversion
    case util.EqualCmd(cmd, "show tms groups"):
        return getFixture("tms/groups", nil, `TMS Groups (emulated)
Name    Diversion   Targets
edge1   flowspec    192.0.2.0/24
edge2   rtbh        198.51.100.0/24

`)
    case strings.HasPrefix(cmd, "show tms group "):
        name := strings.TrimPrefix(cmd, "show tms group ")
        return getFixture("tms/group", map[string]string{"NAME": name}, fmt.Sprintf(`TMS Group %q (emulated)
Diversion: flowspec
Targets:   192.0.2.0/24
Routers:   r1, r2

`, name))
    case util.EqualCmd(cmd, "show diversion status"):
        return getFixture("diversion/status", nil, `Diversion Status (emulated)
Method:    flowspec
State:     programmed
Routes:    12

`)

    // Policy management
    case util.EqualCmd(cmd, "show policies"):
        return getFixture("policies/catalog", policies, formatPolicies())

    // Health & system
    case util.EqualCmd(cmd, "show system"), util.EqualCmd(cmd, "show health"):
        return getFixture("system/health", map[string]string{"UPTIME": time.Since(StartTime()).Round(time.Second).String()}, "System Health (emulated)\n" +
            "Uptime:   " + time.Since(StartTime()).Round(time.Second).String() + "\n" +
            "CPU:      8%\n" +
            "Memory:   62%\n" +
            "Disk:     41%\n\n")
    case util.EqualCmd(cmd, "show processes"):
        return getFixture("system/processes", nil, `Processes (emulated)
PID   Name       CPU  Mem
100   tmsd       3%   120MB
101   flowspecd  1%   60MB

`)
    case util.EqualCmd(cmd, "show license"):
        return getFixture("system/license", nil, `License (emulated)
Serial:   SL-FAKE-1234
Features: mitigation, flowspec, rtbh
Expires:  never

`)

    // BGP / Flowspec
    case util.EqualCmd(cmd, "show bgp summary"):
        return getFixture("bgp/summary", nil, `BGP Summary (emulated)
Neighbor        AS   MsgRcvd  MsgSent  State
203.0.113.1     65001  1024     980      Established

`)
    case util.EqualCmd(cmd, "show bgp neighbors"):
        return getFixture("bgp/neighbors", map[string]string{"UPTIME": time.Since(StartTime()).Round(time.Second).String()}, `BGP Neighbors (emulated)
Neighbor: 203.0.113.1
AS:       65001
State:    Established
Uptime:   {{.UPTIME}}

`)
    case util.EqualCmd(cmd, "show flowspec status"):
        return getFixture("flowspec/status", nil, `FlowSpec Status (emulated)
Sessions:  1
Installed: 12 rules
State:     active

`)

    // Usability
    case util.EqualCmd(cmd, "show clock"):
        return getFixture("system/clock", map[string]string{"NOW": time.Now().UTC().Format(time.RFC3339)}, "Clock (emulated)\n{{.NOW}} UTC\n\n")
    case util.EqualCmd(cmd, "show configuration differences"):
        return getFixture("config/differences", nil, "No configuration differences (emulated).\n\n")

    default:
        return fmt.Sprintf("%% Unrecognized 'show' command: %q\nTry: show version | show running-config | show mitigation | show interfaces\n\n", cmd)
    }
}
