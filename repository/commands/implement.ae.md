---
description: Execute the implementation plan
---

# Implement

Execute the current implementation plan. Follow the plan step by step, validate as you go, keep the repo green.

## Inputs

> [!IMPORTANT]
> Resolve the change set by running `ls changes/ | grep -v archive/`. If exactly one directory exists, use it. Only prompt the user when multiple change sets are present.

## Instructions

Load `spec-driven-development` and `research` skills. Read state.md and tasks.md. Apply state entry check.

Determine lane and load plan: full lane reads tasks.md to find current task ([o] or first [ ]), then loads corresponding plan from plans/; vibe/bug lane reads plan.md (single combined plan).

Execute the plan step by step. Validate after significant changes, keep repo green, document deviations. Use research skill for unexpected code structure, dependency questions, or integration uncertainty. Don't guess when you can research.

Handle issues: minor adjustments—proceed and document deviation; major issues—stop and discuss with user; spec issues (full lane)—flag for reconciliation.

After implementation, run validation steps from plan, verify acceptance criteria, ensure tests pass.

**Full lane completion**: When user approves, mark current task ([o]) as complete ([x]) in tasks.md. If tasks remain, update state.md: `## Phase Status: complete`, clear `## Notes`, suggest `/sdd/plan <name>`. If no tasks remain, suggest `/sdd/reconcile <name>`.

**Vibe/bug lane completion**: Discuss next steps. If throwing away—done, no state update. If keeping—update state: `## Phase Status: complete`, suggest `/sdd/reconcile <name>`. Reconcile is optional for vibe/bug.

## Examples

**Full lane implementing a task:**

```text
Input: None (change: "password-reset")
Output: "Loading plan 01.md to implement password validator changes."
       Follows steps: update validator.ts, add reset logic, update tests.
       Validation: All tests pass.
       User: "Looks good."
       Output: "Marked task 1 complete. Three tasks remaining—suggest /sdd/plan for next task."
```

**Vibe lane quick fix:**

```text
Input: "bug-fix" (user has context with plan.md)
Output: "Following plan.md to patch router/routes.ts."
       Implementation complete, tests pass.
       User: "Great, keep this work."
       Output: "Updated state complete. Suggest /sdd/reconcile if you want to capture specs."
