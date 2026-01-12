# Agent Extensions

BrassworksAI's toolkit for spec-driven development. Drop-in agents and commands for OpenCode and Augment.

## Why SDD

SDD lets you work at the capability level, not the ticket level. Describe what you want the system to do—the agents figure out how to break it down, what order to build it, which files to touch. You're not managing a backlog of micro-tasks. You're defining outcomes and validating that they landed correctly.

This is vibe coding's ambition with enterprise rigor. Think big, stay traceable.

## Installation

### Prerequisites

Install at least one of:
- [OpenCode](https://opencode.ai/docs/#install)
- [Augment CLI](https://www.augmentcode.com/)

### macOS / Linux

```sh
git clone git@github.com:BrassworksAI/agent-extensions.git && cd agent-extensions
./install.sh
```

Global installs use symlinks to the repo; local installs copy files into the repo. macOS/Linux only.

### Windows

```powershell
git clone git@github.com:BrassworksAI/agent-extensions.git; cd agent-extensions
.\install.ps1
```

Global installs use symlinks to the repo; local installs copy files into the repo.

### Install Locations

| Tool | Global | Local |
|------|--------|-------|
| OpenCode | `~/.config/opencode/` | `.opencode/` |
| Augment | `~/.augment/` | `.augment/` |

### Uninstall

```sh
# macOS / Linux
curl -fsSL https://raw.githubusercontent.com/BrassworksAI/agent-extensions/main/uninstall.sh | sh
```

```powershell
# Windows
irm https://raw.githubusercontent.com/BrassworksAI/agent-extensions/main/uninstall.ps1 | iex
```

## Available Commands

### Workflow

| Command | Description |
|---------|-------------|
| `/sdd/plan <goal>` | Generate plan for a new feature |
| `/sdd/specs <goal>` | Draft spec for a new feature |
| `/sdd/proposal <goal>` | Create proposal for a feature |
| `/sdd/implement <goal>` | Implement spec |
| `/sdd/reconcile <goal>` | Validate implementation vs spec |
| `/sdd/finish <goal>` | Finalize implementation |
| `/sdd/continue` | Continue last workflow stage |

### Discovery

| Command | Description |
|---------|-------------|
| `/sdd/brainstorm <goal>` | Generate high-level options |
| `/sdd/explain <subject>` | Explain system/flow |
| `/sdd/discovery <goal>` | Explore requirements and constraints |
| `/sdd/tools/critique [target]` | Critique proposal, specs, or plan |
| `/sdd/tools/scenario-test [name]` | Roleplay as user to stress-test design |
| `/sdd/tools/taxonomy-map <name>` | Map intents to spec paths |
| `/sdd/tools/prime-specs <name>` | Load specs into context |

### Meta

| Command | Description |
|---------|-------------|
| `/create/agent <spec>` | Create a new agent |
| `/create/command <spec>` | Create a new command |

## Development Install

Run the install script from your clone:

```sh
git clone git@github.com:BrassworksAI/agent-extensions.git && cd agent-extensions
./install.sh
```

## Repository Structure

```
agent-extensions/
├── opencode/           # OpenCode extensions
│   ├── agent/          # Agent definitions
│   ├── command/        # Command definitions
│   └── skill/          # Skill definitions
├── augment/            # Augment extensions
│   ├── agents/         # Agent definitions
│   ├── commands/       # Command definitions
│   └── skills/         # Skill files
├── install.sh          # macOS/Linux installer (global symlinks, local copies)
├── install.ps1         # Windows installer (global symlinks, local copies)
├── uninstall.sh        # macOS/Linux uninstaller
└── uninstall.ps1       # Windows uninstaller
```

## Learn More

- [THE_CASE_FOR_SPECS.md](THE_CASE_FOR_SPECS.md)
