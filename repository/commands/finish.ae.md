---
description: Close a change set and sync its specs to canonical
---

### Required Skills (Must Load)

- `spec-driven-development` (state management, completion workflow)

# Finish Change Set

Close the change set and sync specs to canonical.

## Inputs

> [!IMPORTANT]
> Resolve the change set by running `ls changes/ | grep -v archive/`. If exactly one directory exists, use it. Only prompt the user when multiple change sets are present.

## Instructions

1. **Verify prerequisites**: Load `spec-driven-development` skill and read state. Confirm phase is `reconcile` and status is `complete` per skill guidelines.

2. **Sync specs**: Use the CLI merge command:
   - Inspect canonical spec inventory first with `ae sdd spec list` (no flag by default)
   - Run a dry run preview first: `ae sdd spec merge --dry-run [<change-name>] [--specs-dir <dir>]`
   - Present dry run results (created/modified/skipped) and any blockers
   - Wait for user approval
   - Apply merge: `ae sdd spec merge [<change-name>] [--specs-dir <dir>]`

3. **Update state**: Mark phase and status as complete per skill guidelines.

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
