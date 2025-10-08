# Mock Arbor TMS SSH Emulator

(c) 2025 Sam Caldwell <mail@samcaldwell.net>

## Latest Version: 
[![Version](https://img.shields.io/github/v/tag/asymmetric-effort/mock-arbor?sort=semver)](https://asymmetric-effort/mock-arbor/tags)

Minimal SSH server that emulates an Arbor TMS-style CLI for testing and tooling. Supports interactive shell and 
single-command exec with a small set of static responses. Not intended for production use.

## Features
- SSH server with password auth (no pubkey for simplicity)
- Interactive prompt and exec mode
- Emulated commands: 
  - `help`, `show version`, `show running-config|run`, `show mitigation|mitigation status`, `show interfaces|int`, 
  - `show logging`, `show policy <name>`, `mitigation start|stop|apply <policy>`, `configure terminal|conf t`, 
  - `write memory|wr mem`, `clear counters`, `exit|quit|logout`

## Requirements
- Go 1.25+

### Build and Run
- `make build` — builds `bin/tms_ssh_emulator`
- `bin/tms_ssh_emulator` — starts the server

### Defaults and Env Vars
- `TMS_SSH_ADDR` (default `:2222`)
- `TMS_USERNAME` (default `tms`)
- `TMS_PASSWORD` (default `tms123!`)
- `TMS_HOSTKEY_FILE` (default `hostkey.pem`)

## Quick Start
- Interactive: `ssh -p 2222 tms@localhost` (password: `tms123!`)
- Exec mode: `ssh -p 2222 tms@localhost show version`

## Development
- `make test` — race + coverage
- `make cover` — prints summary and writes `out/coverage.html`
- `make lint` — runs `go vet`
- `make clean` — removes build and coverage artifacts

## Tagging (Semantic Versioning)
- `make tag/patch` — bumps `vX.Y.Z` to `vX.Y.(Z+1)`
- `make tag/minor` — bumps `vX.Y.Z` to `vX.(Y+1).0`
- `make tag/major` — bumps `vX.Y.Z` to `v(X+1).0.0`
- Behavior:
  - Finds the highest tag matching `vMAJOR.MINOR.PATCH` (semantic version format).
  - If no tags exist, it creates and pushes `v0.0.0` first, then creates the bumped tag.
  - Creates annotated tags and pushes them to `origin`.

## Notes
- If `TMS_HOSTKEY_FILE` does not exist, an RSA key is generated; best-effort persistence is attempted.
