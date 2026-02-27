---
name: spec-driven-development
description: Spec-Driven Development (SDD) workflow and state management. Use when guiding or executing SDD flows; selecting lanes (full/vibe/bug), managing artifacts, running task progression commands, and keeping change state current through ae sdd CLI.
---

# Spec-Driven Development

Use this skill as the source of truth for running SDD with `ae sdd` commands. Load it before `proposal`, `bug`, or `vibe` work so change-set resolution and command flow are consistent.

## Core Rules

- Use CLI commands for state progression whenever possible.
- Never advance phases automatically from any workflow command.
- Phase advancement is user-driven via `ae-sdd-next` only.
- If the user approves work, suggest `ae-sdd-next` instead of running phase progression commands.
- Do not use custom init scaffolding; initialize with `ae sdd init`.
- Keep one task `in_progress` at a time in full lane.
- Keep `notes` and `pending` current so work can resume in a new chat.

## Lane Selection

- **full**: New capabilities, non-trivial behavior changes, architecture work.
- **vibe**: Exploration and prototyping that should still be tracked.
- **bug**: Defect fixes against intended or specified behavior.

If a bug request is actually a behavior change, switch to `full`.

## Command Quick Start

1. Resolve active change set (or create one):

```bash
ae sdd init <name> --lane <full|vibe|bug>
```

1. Check status before doing work:

```bash
ae sdd status [name]
```

1. Move to the right phase only through `ae-sdd-next`:

```bash
ae-sdd-next [name]
```

## Full Command Reference

### Change Set and Phase Commands

| Command | Purpose |
|---|---|
| `ae sdd init <name> --lane <full|vibe|bug>` | Create a new change set and initial state |
| `ae sdd status [name]` | Show current lane, phase, tasks, notes, and pending |
| `ae sdd config init` | Initialize `.ae-config.json` for canonical spec root |
| `ae-sdd-next [name]` | Verify current-phase artifacts, then complete and advance |
| `ae sdd phase complete [--next] [name]` | Underlying CLI transition used by `ae-sdd-next` |
| `ae sdd phase set <phase> [name]` | Manual override for explicit user-directed corrections |
| `ae sdd phase next [name]` | Underlying CLI step used after completion |

### Task Commands (Full Lane)

| Command | Purpose |
|---|---|
| `ae sdd task list [name]` | Show ordered task list |
| `ae sdd task current [name]` | Show current in-progress task |
| `ae sdd task next [name]` | Show next pending task |
| `ae sdd task start [name]` | Start next pending task (or named task when supported) |
| `ae sdd task complete [name]` | Complete current in-progress task |
| `ae sdd task complete --next [name]` | Complete current task and immediately start next |

### Notes and Pending Commands

| Command | Purpose |
|---|---|
| `ae sdd notes set "content" [name]` | Update resume context and decisions |
| `ae sdd pending add "item" [name]` | Track unresolved blockers |
| `ae sdd pending clear <index> [name]` | Remove resolved blocker |

## Phase Flows

Full lane:

```text
proposal -> specs -> discovery -> tasks -> plan -> implement -> reconcile -> finish
```

Vibe lane:

```text
context -> plan -> implement -> [reconcile -> finish]
```

Bug lane:

```text
triage -> plan -> implement -> [reconcile -> finish]
```

For vibe and bug lanes, `reconcile` and `finish` are optional unless specs must be updated.

## Implement-Phase Rules (Full Lane)

- Start implementation tasks with `ae sdd task start`.
- Keep task order in `tasks.toml`; order is execution priority.
- Complete work using `ae sdd task complete` or `ae sdd task complete --next`.
- `ae sdd phase next` from `implement` is guarded:
  - blocked if any task is currently `in_progress`
  - loops back to `plan` when tasks remain incomplete
  - advances beyond `implement` only when all tasks are complete

## Phase Transition Guardrails

- Never auto-transition phases inside `ae-sdd-proposal`, `ae-sdd-specs`, `ae-sdd-tasks`, `ae-sdd-plan`, `ae-sdd-implement`, `ae-sdd-reconcile`, `ae-sdd-vibe`, `ae-sdd-bug`, or other workflow commands.
- Only `ae-sdd-next` performs phase progression during normal workflow.
- `ae-sdd-next` must verify phase artifacts before running `ae sdd phase complete --next`.
- If artifacts are missing, stop and instruct the user how to produce them.
- `ae sdd phase set <phase>` is reserved for explicit user-directed corrections.

## Artifact Expectations

- `changes/<name>/state.toml`: lane, phase, notes, pending
- `changes/<name>/tasks.toml`: ordered tasks and status (full lane)
- `changes/<name>/proposal.md`: proposal (full lane)
- `changes/<name>/context.md`: exploratory context (vibe/bug as needed)
- `changes/<name>/specs/**/*.md`: specs
- `changes/<name>/plans/*.md`: implementation plans

## Session Playbook

At the start of any SDD command session:

1. Run `ae sdd status [name]`.
2. Confirm lane and phase.
3. For full lane, inspect tasks with `ae sdd task list` and `ae sdd task current`.
4. Execute the phase-appropriate command.
5. Update notes and pending items before ending session.

## Quality Gates

- Use one active change set per thread of work.
- Do not manually invent state structure in `state.toml`.
- Keep pending list strictly unresolved items.
- Keep notes concise and resume-oriented.
- In full lane, do not leave multiple tasks `in_progress`.
