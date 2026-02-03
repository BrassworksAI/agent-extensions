# Agent Extensions

Go CLI that installs commands, skills, and hooks for AI coding agents.

## Quick Reference

**Default to using `mise run <task>` for all operations.** Only use custom commands if a task is not available.

| Task | Command | Description |
|------|---------|-------------|
| build | `mise run build` | Build the ae CLI binary (includes embed) |
| embed | `mise run embed` | Prepare embedded content (tools.yaml + repository/) |
| dev | `mise run dev` | Build and run ae with args |
| test | `mise run test` | Run all tests |
| lint | `mise run lint` | Run golangci-lint |
| fmt | `mise run fmt` | Format Go code |
| tidy | `mise run tidy` | Tidy go modules |
| release | `mise run release` | Build release binaries with GoReleaser |
| install | `mise run install` | Install ae to GOBIN |
| clean | `mise run clean` | Remove build artifacts |
