---
name: sdd/finish
description: Close the change set and sync specs
agent: sdd/plan
---

# Required Skills (Must Load)

You MUST load and follow these skills before doing anything else:

- `sdd-state-management`
- `spec-format`

If any required skill content is missing or not available in context, you MUST stop and ask the user to re-run the command or otherwise provide the missing skill content. Do NOT proceed without them.

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

- **`kind: new`**: copy/move the spec content into canonical under `specs/` at the same relative path.
- **`kind: delta`**: merge the delta into the existing canonical spec.
  - Apply `### ADDED / ### MODIFIED / ### REMOVED` buckets (topics under `####`).
  - MODIFIED uses adjacent `Before/After` to locate and update text.

Verify canonical reflects the intended changes.

**Note:** Delta merging will eventually be automated; for now apply merges carefully and review results with the user.

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
