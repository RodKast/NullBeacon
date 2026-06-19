# 🕸️ go-c2

A Command & Control (C2) framework built from scratch in Go, developed as a learning project to explore Go concurrency, networking, and security research concepts.

> ⚠️ **For authorized security research, CTF, and home lab use only. Never deploy against systems you do not own or have explicit written permission to test.**

---

## 🏗️ Architecture

```
Operator
   │
   ▼
Teamserver  ◄──── TCP :8080 ────  Agent(s)
```

| Component | Description |
|---|---|
| `teamserver` | Listens for agent connections, registers agents, queues tasks |
| `agent` | Beacons to the teamserver every 10 seconds with system info |

---

## 📁 Project Structure

```
go-c2/
├── cmd/
│   ├── teamserver/     # Teamserver binary
│   └── agent/          # Agent binary
├── pkg/
│   ├── agent/          # Agent struct and registration logic
│   ├── transport/      # (upcoming) Transport layer
│   └── task/           # (upcoming) Task queue
├── go.mod
└── go.sum
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

## ✅ Features

- [x] TCP listener with concurrent agent handling
- [x] Read timeout protection against hung connections
- [x] Agent check-in with hostname and username
- [x] UUID-based agent registration
- [x] Mutex-protected agent map for concurrent access
- [x] Configurable server address via CLI flag
- [x] Agent beacon loop (every 10 seconds)

---

## Roadmap

- [ ] Stage 5 — Operator CLI (interactive shell, list agents, queue tasks)
- [ ] Stage 6 — Persistence (registry, cron, systemd)
- [ ] Stage 7 — Evasion (sleep jitter, AMSI/ETW stubs)
- [ ] Stage 8 — Packing (AES payload encryption, custom loaders)
- [ ] Stage 9 — Transport hardening (TLS, malleable profiles)

---

## Disclaimer

This project is built for **educational purposes**, **CTF competitions**, and **authorized security research**. The authors are not responsible for misuse.
