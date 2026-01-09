---
description: Close change set and sync specs
argument-hint: <change-set-name>
---

# Finish

Close change set and sync change-set specs to canonical.

## Arguments

- `$ARGUMENTS` - Change set name

## Instructions

> **SDD Process**: Read `.augment/skills/sdd-state-management.md` for state management guidance.

> **Spec Format**: Use guidance from `.augment/skills/spec-format.md` (project-local) or `~/.augment/skills/spec-format.md` (global) for spec structure.

### Setup

1. Read `changes/<name>/state.md`

### Entry Check

Apply state entry check logic from `.augment/skills/sdd-state-management.md`.

Verify prerequisites: Reconciliation complete (phase `reconcile`, status `complete`).

### Stage 1: Sync Preview (Required)

If `changes/<name>/specs/` exists (created/updated during reconcile):

1. Enumerate all spec files under `changes/<name>/specs/`.
2. For each spec file, read its required YAML frontmatter:

```markdown
---
kind: new | delta
---
```

3. Present a **sync plan preview** to the user (before making any changes):
   - Source: `changes/<name>/specs/<path>.md`
   - Kind: `new` or `delta`
   - Target canonical path: `specs/<path>.md`
   - For `kind: delta`: confirm the target canonical spec exists and call out what sections you plan to add/modify/remove.

4. Call out blockers requiring user decisions:
   - Missing target canonical spec for `kind: delta`
   - Ambiguous `Before/After` matches
   - Any uncertainty about intent or boundaries

You MUST WAIT for the user to explicitly approve the sync plan before applying any spec changes.

### Stage 2: Apply Sync (After Approval)

Only after the user explicitly approves the sync plan:

- **`kind: new`**: copy/move spec content into canonical under `specs/` at same relative path.
- **`kind: delta`**: merge delta into existing canonical spec.
  - Apply `### ADDED / ### MODIFIED / ### REMOVED` buckets (topics under `####`).
  - MODIFIED uses adjacent `Before/After` to locate and update text.

Verify canonical reflects intended changes.

**Note:** Delta merging will eventually be automated; for now apply merges carefully and review results with user.

### Update State

Update `changes/<name>/state.md`:

```markdown
## Phase

complete

## Phase Status

complete
```

Add completion timestamp to notes or leave empty.

### Cleanup Options (Separate Approval)

Discuss cleanup preference with user. Approval to sync specs does NOT imply approval to clean up artifacts.

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
