# 🕸️ go-c2

A Command & Control (C2) framework built from scratch in Go, developed as a learning project to explore Go concurrency, networking, and security research concepts.

> ⚠️ **For authorized security research, CTF, and home lab use only. Never deploy against systems you do not own or have explicit written permission to test.**

---

## 🏗️ Architecture

```
Operator CLI
     │
     ▼
Teamserver  ◄──── TCP :8080 ────  Agent(s)
     │
     └── teamserver.log
```

| Component | Description |
|---|---|
| `teamserver` | Listens for agent connections, registers agents, queues tasks |
| `agent` | Beacons to the teamserver every 10 seconds with system info |
| `operator shell` | Interactive CLI for listing agents and queuing tasks |

---

## 📁 Project Structure

```
go-c2/
├── cmd/
│   ├── teamserver/     # Teamserver binary + operator shell
│   └── agent/          # Agent binary
├── pkg/
│   ├── agent/          # Agent struct and registration logic
│   ├── task/           # Task struct and queue logic
│   └── transport/      # (upcoming) Transport layer
├── go.mod
├── go.sum
└── teamserver.log      # Runtime log (gitignored)
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
| `list` | List all connected agents |
| `interact <agentID>` | Enter the agent shell |
| `exit` | Exit the operator shell |

**Inside the agent shell:**

| Command | Description |
|---|---|
| `<any command>` | Queue a task for the agent |
| `back` | Return to the main shell |

**Example session:**
```
Enter command (list, interact <agent_id>, exit): list
Connected agents:
ID: bbaec000-a7da-499e-8c7c-32e52a41c767, Username: chris, Hostname: target, Address: 192.168.1.5:51234

Enter command (list, interact <agent_id>, exit): interact bbaec000-a7da-499e-8c7c-32e52a41c767
[agent:bbaec000]> whoami
Task created with ID: 62c182fa-1837-47f2-b0d4-79b392cabaef
[agent:bbaec000]> back

Enter command (list, interact <agent_id>, exit): exit
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

---

## 🗺️ Roadmap

- [x] Stage 1 — TCP teamserver skeleton
- [x] Stage 2 — Agent check-in with system info
- [x] Stage 3 — Beacon loop
- [x] Stage 4 — Agent registration with UUID and mutex-protected map
- [x] Stage 5 — Operator shell (list agents, queue tasks)
- [ ] Stage 6 — Persistent agent ID + task delivery on beacon
- [ ] Stage 7 — Command execution on agent side
- [ ] Stage 8 — Persistence (registry, cron, systemd)
- [ ] Stage 9 — Evasion (sleep jitter, AMSI/ETW stubs)
- [ ] Stage 10 — Packing (AES payload encryption, custom loaders)
- [ ] Stage 11 — Transport hardening (TLS, malleable profiles)

---

## ⚠️ Disclaimer

This project is built for **educational purposes**, **CTF competitions**, and **authorized security research**. The authors are not responsible for misuse.
