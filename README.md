# [OpenCode](https://opencode.ai) SDD

**Speed without direction is just chaos with velocity.**

AI writes the code now. SDD makes sure it's the right code—capturing intent before implementation so you don't lose track of what's actually in your codebase.

## How It Works

SDD breaks work into phases: proposal, specs, discovery, tasks, planning, implementation, reconciliation. Sounds heavy. In practice, AI agents do most of the work at each phase—you're reviewing and steering, not writing.

Every artifact ends with a feedback section. Don't like what the agent produced? Write your notes, rerun the command, and it revises. No context lost, no starting over. The loop is tight.

The system catches drift before it ships. Discovery flags architecture conflicts before you've written code. Reconciliation at the end validates that what got built actually matches what you intended. If the implementation wandered, you'll know—and you'll have specs to point at.

You stay in control without doing the busywork. The agents handle the decomposition, the planning, the implementation. You handle the "is this actually what we want?" question. That's the job now.

## Think Bigger

Scrum taught us to think small. Break everything into bite-sized tickets because that's all a human can safely hold. Keep changes tiny because big changes are terrifying when three people understand the system.

AI doesn't have that limitation. Through this process, agents can reason about entire features—see patterns that only make sense when you're looking at the whole picture, even as context flows across phases.

SDD lets you work at the capability level, not the ticket level. Describe what you want the system to do—the agents figure out how to break it down, what order to build it, which files to touch. You're not managing a backlog of micro-tasks. You're defining outcomes and validating that they landed correctly.

This is vibe coding's ambition with enterprise rigor. Think big, stay traceable.

## Installation

### macOS / Linux

```sh
curl -fsSL https://raw.githubusercontent.com/BrassworksAI/opencode-sdd/main/install.sh | sh
```

### Windows (PowerShell)

```powershell
irm https://raw.githubusercontent.com/BrassworksAI/opencode-sdd/main/install.ps1 | iex
```

### Uninstall

```sh
# macOS / Linux
curl -fsSL https://raw.githubusercontent.com/BrassworksAI/opencode-sdd/main/uninstall.sh | sh
```

```powershell
# Windows
irm https://raw.githubusercontent.com/BrassworksAI/opencode-sdd/main/uninstall.ps1 | iex
```

## Getting Started

After installation, execute OpenCode and run:

```
/sdd/explain
```

This walks you through the full workflow, what each phase produces, and how to get started. You can also run `/sdd/explain <phase>` for details on any specific phase.

Or jump straight in:

```
/sdd/init my-feature      # Start a new feature
/sdd/bug fix-login        # Fast lane for bug fixes  
/sdd/quick try-something  # Rapid prototyping
```

## Development Install

For contributors who want to edit files while having them active in OpenCode:

```sh
git clone git@github.com:BrassworksAI/opencode-sdd.git && cd opencode-sdd
```

```sh
./dev-install.sh
```

Creates symlinks so edits in `~/.config/opencode` modify the repo directly. macOS/Linux only.

## Learn More

Read [THE_CASE_FOR_SPECS.md](./THE_CASE_FOR_SPECS.md) for the story of how an AI came to understand why it needed the harness.
