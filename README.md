# Mock Arbor TMS SSH Emulator

(c) 2025 Sam Caldwell <mail@samcaldwell.net>

Minimal SSH server that emulates an Arbor TMS-style CLI for testing and tooling. Supports interactive shell and 
single-command exec with a small set of static responses. Not intended for production use.

Features
- SSH server with password auth (no pubkey for simplicity)
- Interactive prompt and exec mode
- Emulated commands: 
  - `help`, `show version`, `show running-config|run`, `show mitigation|mitigation status`, `show interfaces|int`, 
  - `show logging`, `show policy <name>`, `mitigation start|stop|apply <policy>`, `configure terminal|conf t`, 
  - `write memory|wr mem`, `clear counters`, `exit|quit|logout`

Requirements
- Go 1.25+

Build and Run
- `make build` — builds `bin/tms_ssh_emulator`
- `bin/tms_ssh_emulator` — starts the server

Defaults and Env Vars
- `TMS_SSH_ADDR` (default `:2222`)
- `TMS_USERNAME` (default `tms`)
- `TMS_PASSWORD` (default `tms123!`)
- `TMS_HOSTKEY_FILE` (default `hostkey.pem`)

Quick Start
- Interactive: `ssh -p 2222 tms@localhost` (password: `tms123!`)
- Exec mode: `ssh -p 2222 tms@localhost show version`

Development
- `make test` — race + coverage
- `make cover` — prints summary and writes `out/coverage.html`
- `make lint` — runs `go vet`
- `make clean` — removes build and coverage artifacts

Notes
- If `TMS_HOSTKEY_FILE` does not exist, an RSA key is generated; best-effort persistence is attempted.
