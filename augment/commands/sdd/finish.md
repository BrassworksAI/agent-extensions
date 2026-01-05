---
description: Close the change set and sync specs
argument-hint: <change-set-name>
---

# Finish

Close the change set and sync delta specs to canonical.

## Arguments

- `$ARGUMENTS` - Change set name

## Instructions

### Setup

1. Read `changes/<name>/state.md` - verify phase is `finish`
2. Verify prerequisites: Reconciliation complete

### Sync Delta Specs

If `changes/<name>/specs/` exists (created in reconcile phase):

For each delta spec in `changes/<name>/specs/`:

1. **Read the delta spec**
2. **Apply to canonical**: 
   - New capabilities: Create new file in `specs/`
   - Modified requirements: Update existing file in `specs/`
   - Removed requirements: Remove from existing file in `specs/`
3. **Verify sync**: Ensure canonical reflects all changes

### Update State

Update `changes/<name>/state.md`:

```markdown
## Phase

complete

## Completed

<timestamp>
```

### Cleanup Options

Discuss cleanup preference with user:

1. **Keep all artifacts**: Leave `changes/<name>/` intact for history
2. **Archive**: Move to `changes/archive/<name>/`
3. **Remove**: Delete `changes/<name>/` (delta specs already synced)

Only proceed with cleanup after user explicitly chooses an option.

### Summary Report

Provide completion summary:

- What was accomplished
- Files changed
- Specs added/modified/removed (if specs were created)
- Any notes or follow-up items
