---
name: sdd/implement
description: Execute the implementation plan
agent: sdd/build
---

# Required Skills (Must Load)

You MUST load and follow these skills before doing anything else:

- `sdd-state-management`
- `research`

If any required skill content is missing or not available in context, you MUST stop and ask the user to re-run the command or otherwise provide the missing skill content. Do NOT proceed without them.

# Implement

Execute the current implementation plan.

## Arguments

- `$ARGUMENTS` - Change set name

## Instructions

### Setup

!`cat changes/$1/state.md 2>/dev/null || echo "State file not found"`

!`cat changes/$1/tasks.md 2>/dev/null || echo "No tasks found"`

### Entry Check

Apply state entry check logic from `sdd-state-management` skill.

Determine lane and load plan:
- **Full lane**: Read `changes/<name>/tasks.md` only to identify the current task (prefer `[o]`; otherwise first `[ ]`). Then load the corresponding plan from `changes/<name>/plans/`.
- **Vibe/Bug lane**: Read `changes/<name>/plan.md` (single combined plan)

For full lane, the plan file is the source of truth for what to do; `tasks.md` is bookkeeping (current task + completion).

### Implementation Process

Execute the plan step by step:

1. **Follow the plan**: The plan was created for a reason - follow it
2. **Validate as you go**: Run tests/checks after each significant change
3. **Keep the repo green**: Don't leave broken state
4. **Document deviations**: If you must deviate from plan, note why

Update state.md `## Notes` with implementation progress and any deviations.

### Research During Implementation

If you encounter unexpected situations, use the `research` skill:

- **Unexpected code structure**: Investigate to understand it
- **Need to understand a dependency**: Research how it works
- **Unclear integration**: Research before guessing

Don't guess when you can research. But also don't over-research - the plan should have captured the major research needs.

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
1. Review results with the user and confirm validation is green.
2. When the user explicitly approves the implementation, update `changes/<name>/tasks.md`: change the current `[o]` to `[x]`.
3. If any tasks remain `[ ]` after marking the current task complete, update `changes/<name>/state.md`: `## Phase Status: complete`, clear `## Notes`, and suggest `/sdd/plan <name>`.
4. If no tasks remain `[ ]` (i.e., you just completed the last task in the change set), update `changes/<name>/state.md`: `## Phase Status: complete`, clear `## Notes`, and suggest `/sdd/reconcile <name>`.

**Vibe/Bug Lane:**
1. Implementation complete - discuss with user what's next:
   - **Throwing away**: Done - no state update needed
   - **Keeping the work**: When user decides to keep it, update state: `## Phase Status: complete`, clear `## Notes`, suggest `/sdd/reconcile <name>`

Do not log completion in `## Pending` (that section is for unresolved blockers/decisions only).

> **Note**: For vibe/bug lanes, reconcile is optional. If the work is exploratory or a quick fix that doesn't warrant spec updates, stopping here is perfectly valid.
