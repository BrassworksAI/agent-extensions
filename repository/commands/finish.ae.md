---
description: Close a change set and sync its specs to canonical
---

### Required Skills (Must Load)

Load and follow these skills before doing anything else:

- `spec-driven-development`
- `merge-change-specs`

# Finish Change Set

Close the change set and sync specs to canonical.

## Inputs

> [!IMPORTANT]
> Resolve the change set by running `ls changes/ | grep -v archive/`. If exactly one directory exists, use it. Only prompt the user when multiple change sets are present.

- **Change set name**: Resolve by running `ls -1 changes` ignoring `archive/`. If exactly one directory remains, use it as `<change-set-name>`. Only prompt the user when multiple change sets are present.

## Instructions

1. **Verify prerequisites**: Apply state entry check from `spec-driven-development`. Confirm phase is `reconcile` and status is `complete`.

2. **Sync specs**: Use the `merge-change-specs` skill to:
   - Run a dry run preview showing which specs will be created/modified
   - Present any errors or blockers requiring user decisions
   - Wait for user approval
   - Apply the merge to canonical `specs/`

3. **Update state**: Set phase and status to `complete` in `changes/<name>/state.md`.

4. **Cleanup** (separate approval):
   - Keep artifacts intact
   - Archive to `changes/archive/<name>/`
   - Delete `changes/<name>/` (specs synced)

   Ask user to choose.

5. **Report**: Summarize merge results (created/modified/skipped counts), state update, and cleanup action.

## Example

```text
Verified phase=reconcile, status=complete
Dry run: merge 3 new, 2 modified specs
User approved merge
Applied: specs/a.md, specs/b.md, specs/c.md created; specs/d.md modified
State updated to complete
Archived to changes/archive/feature-x/
```

Success: change set closed, specs synced, state updated, cleanup completed per preference.
