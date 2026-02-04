---
description: Show status of SDD change set
---

# Status

Show the status of an SDD change set.

## Required Skills

- `spec-driven-development` (state interpretation, status rendering)

## Inputs

> [!IMPORTANT]
> Resolve the change set by running `ls changes/ | grep -v archive/`. If exactly one directory exists, use it. Only prompt the user when multiple change sets are present.

## Instructions

1. Load the `spec-driven-development` skill for state handling conventions.

2. Read `changes/<change-set-name>/state.toml` and `changes/<change-set-name>/tasks.toml` to get the current state.

3. Report: phase, status, lane, status notes, and next action based on the phase-status mapping from the skill.

4. If in `plan` or `implement` phase, include task progress breakdown using task status values (done, active, pending) from tasks.toml.

5. Output format shows change set metadata, status notes, and the specific command to run next based on current phase and status per skill guidelines.

## Examples

**Single change set exists:**

```text
Input: None (only one change: "login-auth")
Output: Report phase/discovery/status, show complete/in_progress/pending task counts, suggest "/sdd/tasks login-auth"
```

**Multiple change sets exist:**

```text
Input: "What change set would you like status for?" → "login-auth"
Output: Report current state and next action for login-auth
```
