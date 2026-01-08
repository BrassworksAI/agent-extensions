---
description: Execute implementation plan
argument-hint: <change-set-name>
---

# Implement

Execute current implementation plan.

## Arguments

- `$ARGUMENTS` - Change set name

## Instructions

> **SDD Process**: Read `.augment/skills/sdd-state-management.md` for state management guidance.

> **Research**: When needed, delegate to `@librarian` for codebase context. See `.augment/skills/research.md` (project-local) or `~/.augment/skills/research.md` (global) for guidance.

### Setup

1. Read `changes/<name>/state.md`
2. Determine lane and load plan:
   - **Full lane**: Read `tasks.md` only to identify current task (prefer `[o]`; otherwise first `[ ]`). Then load corresponding plan from `plans/`.
   - **Vibe/Bug lane**: Read `plan.md` (single combined plan)

### Entry Check

Apply state entry check logic from `.augment/skills/sdd-state-management.md`.

For full lane, plan file is source of truth for what to do; `tasks.md` is bookkeeping (current task + completion).

### Implementation Process

Execute plan step by step:

1. **Follow plan**: The plan was created for a reason - follow it
2. **Validate as you go**: Run tests/checks after each significant change
3. **Keep repo green**: Don't leave broken state
4. **Document deviations**: If you must deviate from plan, note why

Update state.md `## Notes` with implementation progress and any deviations.

### Research During Implementation

If you encounter unexpected situations, delegate to `@librarian`:

- **Unexpected code structure**: Investigate to understand it
- **Need to understand a dependency**: Research how it works
- **Unclear integration**: Research before guessing

Don't guess when you can research. But also don't over-research - plan should have captured major research needs.

### Handling Issues

If implementation reveals plan problems:

- **Minor adjustments**: Proceed, document deviation in plan
- **Major issues**: Stop, discuss with user, potentially re-plan
- **Spec issues (full lane)**: Flag for reconciliation (don't modify specs during implement)

### Validation

After implementation:

1. Run validation steps from plan
2. Verify acceptance criteria are met
3. Ensure tests pass

### Completion

**Full Lane:**
1. Review results with user and confirm validation is green.
2. When user explicitly approves implementation, update `changes/<name>/tasks.md`: change current `[o]` to `[x]`.
3. If any tasks remain `[ ]` after marking current task complete, update `changes/<name>/state.md`: `## Phase Status: complete`, clear `## Notes`, and suggest `/sdd:plan <name>`.
4. If no tasks remain `[ ]` (i.e., you just completed last task in change set), update `changes/<name>/state.md`: `## Phase Status: complete`, clear `## Notes`, and suggest `/sdd:reconcile <name>`.

**Vibe/Bug Lane:**
1. Implementation complete - discuss with user what's next:
   - **Throwing away**: Done - no state update needed
   - **Keeping work**: When user decides to keep it, update state: `## Phase Status: complete`, clear `## Notes`, suggest `/sdd:reconcile <name>`

Do not log completion in `## Pending` (that section is for unresolved blockers/decisions only).

> **Note**: For vibe/bug lanes, reconcile is optional. If work is exploratory or a quick fix that doesn't warrant spec updates, stopping here is perfectly valid.
