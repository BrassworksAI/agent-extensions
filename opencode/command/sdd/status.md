---
name: sdd/status
description: Show status of injected change set
agent: sdd/plan
---

# Required Skills (Must Load)

You MUST load and follow these skills before doing anything else:

- `sdd-state-management`

If any required skill content is missing or not available in context, you MUST stop and ask the user to re-run the command or otherwise provide the missing skill content. Do NOT proceed without them.

# Status

Show the status of an SDD change set.

## Arguments

- `$ARGUMENTS` - Required: change set name

## Instructions

### Check Arguments

If `$ARGUMENTS` is not provided or empty:
- Stop immediately
- Ask: "What change set would you like me to report on?"
- Do not proceed with any other steps

### Read Injected State Document

The state.md and tasks.md content below is injected from the change set:

!`cat changes/$1/state.md 2>/dev/null || echo "State file not found"`

!`cat changes/$1/tasks.md 2>/dev/null || echo "No tasks found"`

Report:

- Current phase and status
- Lane
- Status notes (from `## Notes`)
- Task progress if in plan/implement phase (e.g., [x] 2, [o] 1, [ ] 5)
- Next suggested action

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
