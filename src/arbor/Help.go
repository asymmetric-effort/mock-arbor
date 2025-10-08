package arbor

// Help - returns the emulator's help menu.
func Help() string {
    return `Available commands (emulated):
  help | ?                       Show this help
  show version                   Show software version
  show running-config | show run Show current configuration
  show mitigation [status]       Show mitigation state
  show mitigation summary        Show mitigation summary
  show mitigation detail <id>    Show mitigation detail
  show mitigation counters [id]  Mitigation counters
  mitigation list                List mitigations
  clear mitigation counters [id] Clear mitigation counters
  show interfaces | show int     Show interfaces summary
  show interfaces brief          Interfaces one-line summary
  show logging                   Show logging settings
  show tms groups                List TMS groups
  show tms group <name>          Show TMS group detail
  show diversion status          Show current diversion status
  show policies                  Show policy catalog (emulated)
  policy <name> [countermeasures ...]  Create/update policy
  no policy <name>               Delete policy
  show system | show health      System health
  show processes                 Process list
  show license                   License info
  show bgp summary               BGP summary
  show bgp neighbors             BGP neighbors
  show flowspec status           Flowspec diversion status
  show configuration differences Show pending configuration diff
  show clock                     Current time
  terminal length <n>            Set paging length (emulated)
  show policy <name>             Show a specific policy
  mitigation start|stop          Control mitigation state
  mitigation apply <policy>      Apply a policy
  services sp backup create full Create a full SP backup (Sightline)
  configure terminal | conf t    Enter config mode (no-op)
  write memory | wr mem          Save configuration (no-op)
  clear counters                 Clear counters (no-op)
  exit | quit | logout           Disconnect
`
}
