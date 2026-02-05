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

Vibe sequence:

```text
context -> plan -> implement -> [reconcile -> finish]
```

Bug sequence:

```text
triage -> plan -> implement -> [reconcile -> finish]
```

Reconcile and finish are optional for vibe/bug lanes unless changes must be reflected in specs.

## Artifacts By Phase (Full Lane)

- `proposal`: `proposal.md` (external doc allowed; state tracks phase)
- `specs`: `changes/<name>/specs/**/*.md` (follow the `spec-format` skill for spec structure)
- `discovery`: `changes/<name>/thoughts/*.md`
- `tasks`: `changes/<name>/tasks.toml`
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
[[task]]
name = "db-models"
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
| `ae sdd pending add "item"` | Add pending item |
| `ae sdd pending clear <index>` | Remove pending item |
| `ae sdd notes set "content"` | Update notes |

## Managing State and Tasks

This section provides the definitive guidance for all state and task operations. All SDD commands must follow these patterns.

### Reading State and Tasks

Always read both files at the start of any SDD command:

1. Read `changes/<name>/state.toml` to get current phase, status, lane, notes, and pending items
2. Read `changes/<name>/tasks.toml` to get task list and current task status
3. Use `ae sdd status [name]` for a styled overview (optional)

### Phase Management

**Set phase to in_progress:**

```bash
ae sdd phase set <phase-name>
```

**Mark phase complete:**

```bash
ae sdd phase set <phase-name>  # Ensures phase is current
# Then update status in state.toml or via completion workflow
```

When a phase finishes, always update `status = "complete"` in state.toml and clear the notes section. Do not auto-advance to the next phase.

### Task Management (Full Lane)

Tasks are edited directly in `changes/<name>/tasks.toml`. Keep tasks in the desired execution order; file order is the source of truth.

```toml
[[task]]
name = "db-models"
title = "Foundation - DB models and migrations"
description = "Add user tables and migrations."
status = "pending"  # pending | in_progress | complete
requirements = [
  "When a new user is created, the system shall persist profile fields."
]
```

Update task status by editing the `status` field.

**Task Status Flow:**

- `pending` → `in_progress` → `complete`
- Only one task should be `in_progress` at a time
- Tasks should be completable in a single session

### Notes Management

Update notes to capture decisions, blockers, and context:

```bash
ae sdd notes set "Your notes here"
```

Keep notes concise but informative enough that a new chat can resume work quickly. Clear notes when marking a phase complete.

### Pending Items

Track blockers and waiting items:

```bash
ae sdd pending add "Waiting for API review"
```

Remove items when resolved:

```bash
ae sdd pending clear 0  # Removes first item
```

Keep only unresolved items. Do not use pending as an approval log.

### Common Workflows

**Starting a new phase:**

1. Read current state.toml
2. Confirm phase transition is valid per Phase Gates
3. `ae sdd phase set <new-phase>`
4. Update notes with phase objectives

**Completing work in a phase:**

1. Complete all phase deliverables
2. For full lane: ensure all tasks are `complete`
3. Clear notes in state.toml
4. Set `status = "complete"` in state.toml
5. Suggest next command per Phase Gates

**Resuming after interruption:**

1. Read state.toml for current phase/status
2. Read tasks.toml for task progress
3. Check notes for context
4. Continue from current state

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

## Required Skills Reference

When using this skill in commands:

- List `spec-driven-development` as a required skill for state/task management
- Follow the Managing State and Tasks section for all file operations
- Never manually edit state.toml structure directly
- Edit tasks.toml directly using the documented format
- Use the documented commands or direct field updates as documented above
