---
name: sdd/reconcile
description: Ensure specs match implementation
agent: sdd/plan
---

<skill>sdd-state-management</skill>
<skill>spec-format</skill>

# Reconcile

Ensure that change set specs match the implementation diff.

**This is a collaborative process** - present your findings to the user and get their input before making changes. Do not make unilateral decisions about what to add, remove, or modify in specs.

## Arguments

- `$ARGUMENTS` - Change set name

## Instructions

### Setup

!`cat changes/$1/state.md 2>/dev/null || echo "State file not found"`

!`cat changes/$1/tasks.md 2>/dev/null || echo "No tasks found"`

### Entry Check

Apply state entry check logic from `sdd-state-management` skill.

### The Process

1. **Get the implementation diff**: What code was actually changed?

2. **If specs/ exists** (`changes/<name>/specs/`):
   - Compare specs to the diff
   - Present findings to the user:
     - Does the diff match what the specs describe?
     - Are there implementation changes not covered by specs?
     - Are there specs describing things not in the diff?
   - Get user input on what to do:
     - Add missing specs for unspecced implementation?
     - Remove specs that don't match diff?
     - Modify existing specs to match implementation?

3. **If specs/ does not exist**:
   - Analyze whether the implementation adds/removes logic worth specifying
   - Present your analysis: what changed and whether it's spec-worthy
   - Ask user: "Should I capture specs for these changes?"
   - If yes: Create `changes/<name>/specs/` and write change-set specs (`kind: new` and/or `kind: delta`)
   - If no: Document that specs were not created (trivial changes)

4. **Document findings** in `changes/<name>/reconciliation.md`

Update state.md `## Notes` with reconciliation findings and decisions.

### Writing Change-Set Specs from Diff

When creating specs from the implementation:

- Analyze what changed and what it means for the system
- Use spec-format skill to write proper change-set specs:
  - Describe added capabilities (positive requirements)
  - Describe removed capabilities (negative requirements)
  - Describe behavioral changes
- Each spec file should cover a logical area of change

### Reconciliation Report

```markdown
# Reconciliation: <name>

## Summary

- Specs updated to match implementation: <yes/no>
- New specs created: <yes/no>
- Unspecced implementation: <none / list items>

## Findings

- Any mismatches found and how they were resolved
- Notes on significant changes

## Next Steps

Proceed to finish
```

### Completion

When they explicitly approve your findings and any spec changes:

1. Update state.md: `## Phase Status: complete`, clear `## Notes`
2. Suggest running `/sdd/finish <name>`

Do not log completion in `## Pending` (that section is for unresolved blockers/decisions only).

**Note**: If change-set specs were created or updated, finish will move `kind: new` specs and merge `kind: delta` specs into canonical.
