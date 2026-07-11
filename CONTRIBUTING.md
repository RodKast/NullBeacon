# Contributing to NullBeacon

Thank you for your interest in contributing to NullBeacon.

## Requirements

- Go 1.26+
- `golangci-lint` — `brew install golangci-lint`
- `govulncheck` — `go install golang.org/x/vuln/cmd/govulncheck@latest`

## Getting Started

```bash
git clone https://github.com/RodKast/NullBeacon.git
cd NullBeacon
go mod tidy
```

## Branch Naming

| Type | Format | Example |
|---|---|---|
| Feature | `feat-<short-description>` | `feat-http-listener` |
| Bug fix | `fix-<short-description>` | `fix-agent-reconnect` |
| Docs | `docs-<short-description>` | `docs-readme-update` |

## Commit Style

We use [Gitmoji](https://gitmoji.dev) commit conventions:

| Emoji | When to use |
|---|---|
| ✨ | New feature |
| 🔧 | Fix or configuration |
| 📝 | Documentation |
| 🐛 | Bug fix |
| ♻️ | Refactor |
| 🔒 | Security fix |

Example:
```
✨ feat: add HTTPS listener with /beacon and /result routes
```

## Before Submitting a PR

1. Run `make fmt` — format your code
2. Run `make lint` — fix all lint errors
3. Run `go build ./...` — confirm it compiles
4. Add `Closes #<issue>` in your PR description

## Code Style

- No comments unless the **why** is non-obvious
- Handle all errors — never ignore with `_`
- OS-specific code goes in `_windows.go` / `_linux.go` files with build tags
- Keep functions small and focused

## Disclaimer

All contributions must be for **authorized security research, CTF, and home lab use only**. Do not contribute techniques specifically designed to evade named commercial security products.
