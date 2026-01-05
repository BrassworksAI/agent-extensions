---
name: sdd/reconcile
description: Ensure specs match implementation
agent: sdd/forge
---

<skill>sdd-state-management</skill>
<skill>spec-format</skill>

# Reconcile

Ensure that delta specs match the implementation diff.

**This is a collaborative process** - present your findings to the user and get their input before making changes. Do not make unilateral decisions about what to add, remove, or modify in specs.

## Arguments

- `$ARGUMENTS` - Change set name

## Instructions

### Setup

1. Read `changes/<name>/state.md` - verify phase is `reconcile`
2. Read `changes/<name>/tasks.md` for context

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
    - Ask user: "Should I create delta specs for these changes?"
    - If yes: Create `changes/<name>/specs/` and write delta specs
    - If no: Document that specs were not created (trivial changes)

4. **Document findings** in `changes/<name>/reconciliation.md`

### Writing Delta Specs from Diff

When creating specs from the implementation:

- Analyze what changed and what it means for the system
- Use spec-format skill to write proper delta specs:
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

Work through reconciliation collaboratively with the user. When they explicitly approve your findings and any spec changes:

1. Log approval in state.md under `## Pending`:
   ```
   None - Reconciliation approved: [brief summary of findings]
   ```
2. Update state.md phase to `finish`
3. Suggest running `/sdd/finish <name>`

Don't advance until the user clearly signals approval. Questions, feedback, or acknowledgments don't count as approval.

**Note**: If delta specs were created or updated, they will be synced to canonical specs during finish.
