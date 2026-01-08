---
name: sdd/status
description: Show status of change sets
agent: sdd/plan
---

<skill>sdd-state-management</skill>

# Status

Show the status of SDD change sets.

## Arguments

- `$ARGUMENTS` - Optional: specific change set name, or "all"

## Instructions

### Specific Change Set

If name provided:

!`cat changes/$1/state.md 2>/dev/null || echo "State file not found"`

!`cat changes/$1/tasks.md 2>/dev/null || echo "No tasks found"`

Report:

- Current phase and status
- Lane
- Status notes (from `## Notes`)
- Task progress if in plan/implement phase (e.g., [x] 2, [o] 1, [ ] 5)
- Next suggested action

### All Change Sets

If "all" or no arguments:

1. Find all `changes/*/state.md`
2. For each, report summary:
    - Name
    - Phase and status
    - Lane
    - Brief status (including task completion ratio if in tasks/plan/implement phase)

### Output Format

```markdown
## Change Set: <name>

**Phase:** <phase> (<status>)
**Lane:** <lane>

### Status
<## Notes content>

### Next Action
<suggested command>
```

### Next Action Logic

| Phase | Status | Next Action |
|-------|--------|-------------|
| `ideation` | `in_progress` | Continue brainstorming |
| `ideation` | `complete` | `/sdd/proposal <name>` |
| `proposal` | `in_progress` | Continue refining proposal |
| `proposal` | `complete` | `/sdd/specs <name>` |
| `specs` | `in_progress` | Continue writing specs |
| `specs` | `complete` | `/sdd/discovery <name>` |
| `discovery` | `in_progress` | Continue architecture review |
| `discovery` | `complete` | `/sdd/tasks <name>` |
| `tasks` | `in_progress` | Continue task breakdown |
| `tasks` | `complete` | `/sdd/plan <name>` |
| `plan` | `in_progress` | Continue planning current task |
| `plan` | `complete` | `/sdd/implement <name>` |
| `implement` | `in_progress` | Continue implementation |
| `implement` | `complete` | `/sdd/reconcile <name>` |
| `reconcile` | `in_progress` | Continue reconciliation |
| `reconcile` | `complete` | `/sdd/finish <name>` |
| `finish` | `complete` | Change set complete - nothing to do |

### Task Progress (Plan/Implement Phases Only)

If phase is `plan` or `implement`, also include:

```markdown
### Tasks: [x] <done>, [o] <active>, [ ] <pending>
```
