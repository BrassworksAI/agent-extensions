---
name: spec-driven-development
description: Spec-Driven Development (SDD) workflow and state management. Use when guiding or executing SDD flows: selecting lanes (full/vibe/bug), moving through phases, managing change set artifacts (specs, tasks, plans, reconcile), updating SDD state and tasks TOML files, or using ae sdd CLI commands.
---

# Spec-Driven Development

Provide an end-to-end guide for running SDD change sets, including lane selection, phase flow, artifacts, and state tracking. Use this skill to orient a new chat to the current SDD state and to keep the workflow consistent.

## Quick Start

1. Determine the lane: `full` (spec-first), `vibe` (fast exploration), or `bug` (fix).
2. Read `changes/<name>/state.toml` and `changes/<name>/tasks.toml` if they exist.
3. Confirm current phase and next intended action before doing work.
4. Update state and tasks as work progresses.

## Lane Selection

- **Full**: New features, architectural changes, or anything that needs formal specs.
- **Vibe**: Exploration and rapid prototyping.
- **Bug**: Defect fixes. If the request is a behavioral change, redirect to full lane.

## Phases And Flows

Full lane sequence:

```text
proposal -> specs -> discovery -> tasks -> plan -> implement -> reconcile -> finish
```

Vibe and bug sequences:

```text
context -> plan -> implement -> [reconcile -> finish]
```

Reconcile and finish are optional for vibe/bug lanes unless changes must be reflected in specs.

## Artifacts By Phase (Full Lane)

- `proposal`: `proposal.md` (external doc allowed; state tracks phase)
- `specs`: `changes/<name>/specs/**/*.md` (follow the `spec-format` skill for spec structure)
- `discovery`: `changes/<name>/thoughts/*.md`
- `tasks`: `changes/<name>/tasks.md`
- `plan`: `changes/<name>/plans/NN.md`
- `implement`: code changes
- `reconcile`: reconciliation report
- `finish`: merge or move specs to canonical

## State And Tasks Files

State is tracked in `changes/<name>/state.toml`:

```toml
[change]
name = "my-feature"
lane = "full"  # full | vibe | bug
created_at = 2026-01-27T10:00:00Z

[phase]
current = "plan"
status = "in_progress"  # in_progress | complete | blocked

[pending]
items = [
  "Waiting on API design decision"
]

[notes]
content = """
Working through auth endpoint design.
Decided to use JWT over sessions.
"""
```

Tasks live in `changes/<name>/tasks.toml`:

```toml
[task.db-models]
title = "Foundation - DB models and migrations"
description = "Add user tables and migrations."
status = "complete"  # pending | in_progress | complete
requirements = [
  "When a new user is created, the system shall persist profile fields."
]
```

## CLI Commands

| Command | Description |
|---------|-------------|
| `ae sdd init <name> --lane <full|vibe|bug>` | Create change set with state.toml |
| `ae sdd status [name]` | Show styled status UI |
| `ae sdd phase next` | Advance to next phase |
| `ae sdd phase set <phase>` | Set current phase |
| `ae sdd task list` | List all tasks |
| `ae sdd task add <short-name>` | Add new task |
| `ae sdd task start <short-name>` | Start a task |
| `ae sdd task complete <short-name>` | Complete a task |
| `ae sdd pending add "item"` | Add pending item |
| `ae sdd pending clear <index>` | Remove pending item |
| `ae sdd notes set "content"` | Update notes |

## Phase Gates (Full Lane)

| From | To | Gate |
|------|----|------|
| proposal | specs | Proposal approved |
| specs | discovery | Change-set specs written (`kind: new|delta`) |
| discovery | tasks | Architecture review complete |
| tasks | plan | Tasks defined |
| plan | implement | Plan approved |
| implement | reconcile | All tasks complete |
| reconcile | finish | Implementation matches specs |

## Phase Management Rules

Use flexible progression. Warn on big jumps but do not block the user.

1. Natural progression: set phase `in_progress` and continue.
2. Jump forward: confirm the skip, then update phase to target.
3. Jump backward: confirm reset to `in_progress`.
4. Re-enter complete phase: confirm, then set to `in_progress`.

When a phase finishes, set `status = complete` but do not auto-advance.

## Pending Semantics

- Keep only unresolved items.
- Remove resolved items (do not strike through).
- Leave empty when nothing is pending.
- Do not use pending as an approval log.

## Notes

Keep `notes.content` updated with decisions, progress, and blockers so a new chat can resume quickly.
