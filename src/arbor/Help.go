package arbor

// Help - returns the emulator's help menu.
func Help() string {
	return `Available commands (emulated):
  help | ?                       Show this help
  show version                   Show software version
  show running-config | show run Show current configuration
  show mitigation [status]       Show mitigation state
  show interfaces | show int     Show interfaces summary
  show logging                   Show logging settings
  show policy <name>             Show a specific policy
  mitigation start|stop          Control mitigation state
  mitigation apply <policy>      Apply a policy
  configure terminal | conf t    Enter config mode (no-op)
  write memory | wr mem          Save configuration (no-op)
  clear counters                 Clear counters (no-op)
  exit | quit | logout           Disconnect
`
}
