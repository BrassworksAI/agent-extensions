---
description: Execute the implementation plan
---

# Implement

Execute the current implementation plan. Follow the plan step by step, validate as you go, keep the repo green.

## Required Skills

- `spec-driven-development` (state/task management, lane detection)
- `research`

## Inputs

> [!IMPORTANT]
> Resolve the change set by running `ls changes/ | grep -v archive/`. If exactly one directory exists, use it. Only prompt the user when multiple change sets are present.

## Instructions

1. Load `spec-driven-development` skill and read state from `changes/<name>/state.toml`. Apply state entry check per skill guidelines.

2. Read tasks from `changes/<name>/tasks.toml` if full lane.

3. Determine lane and load plan:
   - **Full lane**: Identify current task (in_progress status). Read corresponding plan from `changes/<name>/plans/`
   - **Vibe/Bug lane**: Read `changes/<name>/plan.md` (single combined plan)

4. Execute the plan step by step:
   - Follow steps exactly as written
   - Validate after significant changes
   - Keep repo green
   - Use `research` skill for unexpected code structure or integration questions
   - Document any deviations from plan

5. Handle issues:
   - **Minor adjustments**: Proceed and document deviation
   - **Major issues**: Stop and discuss with user
   - **Spec issues** (full lane): Flag for reconciliation

6. Run validation steps from plan, verify acceptance criteria, ensure tests pass.

7. **Full lane completion**:
   - Mark current task complete: `ae sdd task complete <short-name>`
   - If tasks remain, suggest `/sdd/plan <name>` for next task
   - If no tasks remain, update phase status per skill and suggest `/sdd/reconcile <name>`

8. **Vibe/Bug lane completion**:
   - Discuss next steps with user
   - If keeping work: update phase status per skill and suggest `/sdd/reconcile <name>` (optional)
   - If throwing away: done, no state update needed

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
```
