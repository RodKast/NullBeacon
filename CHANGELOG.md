# Changelog

All notable changes to NullBeacon are documented here.

Format follows [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).
Versioning follows [Semantic Versioning](https://semver.org/).

---

## [Unreleased]

### In Progress
- Stage 15 — Syscall obfuscation
- Stage 16 — In-memory execution

---

## [v0.1.2] — 2026-06-22

### Added
- ARM64 binary in release pipeline (`nullbeacon-linux-arm64`)
- CI/CD: lint (golangci-lint), vulnerability scan (govulncheck), cross-platform build checks
- Malleable C2 profiles — configurable beacon URL, User-Agent, interval via `profile.json`
- HTTP beacon agent — agents now beacon over HTTPS using `beaconHTTP` and `sendResult`
- LOLBins — PowerShell encoded command execution, certutil file download (Windows)
- Process injection — VirtualAllocEx, WriteProcessMemory, CreateRemoteThread (Windows)
- Makefile — `build`, `build-agent`, `lint`, `fmt`, `clean`, `release` targets
- `SECURITY.md`, `CONTRIBUTING.md`, `CHANGELOG.md`

### Fixed
- Unchecked error returns flagged by golangci-lint
- Deprecated `rand.Seed` replaced with `rand.New(rand.NewSource(...))`
- Updated to Go 1.26.5 to fix `crypto/tls` vulnerability (GO-2026-5856)

---

## [v0.1.1] — 2026-06-20

### Fixed
- Added `permissions: contents: write` to release workflow to fix 403 error on release creation

---

## [v0.1.0] — 2026-06-20

### Added
- Initial public release
- TLS-only teamserver with dynamic listener management
- HTTPS listener with `/beacon` and `/result` routes
- Cross-platform agent generation (Linux amd64/arm64)
- Windows evasion: AMSI patch, ETW stub
- AES-256-GCM + XOR payload encryption with per-agent key injection
- Linux cron and Windows registry persistence
- readline operator shell with colored output
- GitHub Actions CI pipeline
- Install script (`install.sh`) with auto-architecture detection
- Uninstall command (`nullbeacon --uninstall`)
