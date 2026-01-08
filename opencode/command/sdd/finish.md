---
name: sdd/finish
description: Close the change set and sync specs
agent: sdd/plan
---

<skill>sdd-state-management</skill>
<skill>spec-format</skill>

# Finish

Close the change set and sync change-set specs to canonical.

## Arguments

- `$ARGUMENTS` - Change set name

## Instructions

### Setup

!`cat changes/$1/state.md 2>/dev/null || echo "State file not found"`

### Entry Check

Apply state entry check logic from `sdd-state-management` skill.

Verify prerequisites: Reconciliation complete (phase `reconcile`, status `complete`).

### Sync Change-Set Specs

If `changes/<name>/specs/` exists (created/updated during reconcile):

1. Enumerate all spec files.
2. For each spec file, read its required YAML frontmatter:

```markdown
---
kind: new | delta
---
```

3. Sync behavior:

- **`kind: new`**: copy/move the spec content into canonical under `specs/` at the same relative path.
- **`kind: delta`**: merge the delta into the existing canonical spec.
  - Apply `### ADDED / ### MODIFIED / ### REMOVED` buckets (topics under `####`).
  - MODIFIED uses adjacent `Before/After` to locate and update text.

4. Verify canonical reflects the intended changes.

**Note:** delta merging will eventually be automated; for now apply merges carefully and review with the user.

### Update State

Update `changes/<name>/state.md`:

```markdown
## Phase

complete

## Phase Status

complete
```

Add completion timestamp to notes or leave empty.

### Cleanup Options

Discuss cleanup preference with user:

1. **Keep all artifacts**: Leave `changes/<name>/` intact for history
2. **Archive**: Move to `changes/archive/<name>/`
3. **Remove**: Delete `changes/<name>/` (change-set specs already synced)

Only proceed with cleanup after user explicitly chooses an option.

### Summary Report

Provide completion summary:

- What was accomplished
- Files changed
- Specs added/modified/removed (if specs were created)
- Any notes or follow-up items
