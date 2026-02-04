---
description: Show status of SDD change set
---

# Status

Show the status of an SDD change set.

## Inputs

> [!IMPORTANT]
> Resolve the change set by running `ls changes/ | grep -v archive/`. If exactly one directory exists, use it. Only prompt the user when multiple change sets are present.

## Instructions

Load the `spec-driven-development` skill for state handling conventions.

Read `changes/<change-set-name>/state.md` and `changes/<change-set-name>/tasks.md` to get the current state. Report: phase, status, lane, status notes, and next action based on the phase-status mapping. If in `plan` or `implement` phase, include task progress breakdown (done, active, pending).

Output format shows change set metadata, status notes, and the specific command to run next based on current phase and status.

## Examples

**Single change set exists:**

```text
Input: None (only one change: "login-auth")
Output: Report phase/discovery/status, show [x] [o] [ ] task counts, suggest "/sdd/tasks login-auth"
```

**Multiple change sets exist:**

```text
Input: "What change set would you like status for?" → "login-auth"
Output: Report current state and next action for login-auth
```
