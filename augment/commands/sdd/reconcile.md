---
description: Ensure specs match implementation
argument-hint: <change-set-name>
---

# Reconcile

Ensure that change set specs match implementation diff.

**This is a collaborative process** - present your findings to user and get their input before making changes. Do not make unilateral decisions about what to add, remove, or modify in specs.

## Arguments

- `$ARGUMENTS` - Change set name

## Instructions

> **SDD Process**: Read `.augment/skills/sdd-state-management.md` for state management guidance.

> **Spec Format**: Use guidance from `.augment/skills/spec-format.md` (project-local) or `~/.augment/skills/spec-format.md` (global) for spec structure.

### Setup

1. Read `changes/<name>/state.md`
2. Read `changes/<name>/tasks.md`

### Entry Check

Apply state entry check logic from `.augment/skills/sdd-state-management.md`.

### The Process

1. **Get implementation diff**: What code was actually changed?

2. **If specs/ exists** (`changes/<name>/specs/`):
   - Compare specs to diff
   - Present findings to user:
     - Does diff match what specs describe?
     - Are there implementation changes not covered by specs?
     - Are there specs describing things not in diff?
   - Get user input on what to do:
     - Add missing specs for unspecced implementation?
     - Remove specs that don't match diff?
     - Modify existing specs to match implementation?

3. **If specs/ does not exist**:
   - Analyze whether implementation adds/removes logic worth specifying
   - Present your analysis: what changed and whether it's spec-worthy
   - Ask user: "Should I capture specs for these changes?"
   - If yes: Create `changes/<name>/specs/` and write change-set specs (`kind: new` and/or `kind: delta`)
   - If no: Document that specs were not created (trivial changes)

4. **Document findings** in `changes/<name>/reconciliation.md`

Update state.md `## Notes` with reconciliation findings and decisions.

### Writing Change-Set Specs from Diff

When creating specs from implementation:

- Analyze what changed and what it means for system
- Use spec format guidance to write proper change-set specs:
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
2. Suggest running `/sdd:finish <name>`

Do not log completion in `## Pending` (that section is for unresolved blockers/decisions only).

**Note**: If change-set specs were created or updated, finish will move `kind: new` specs and merge `kind: delta` specs into canonical.
