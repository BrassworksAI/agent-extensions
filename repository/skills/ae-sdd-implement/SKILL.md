---
name: ae-sdd-implement
description: Execute the implementation plan
---

# Implement

Execute the current implementation plan. Follow the plan step by step, validate as you go, keep the repo green.

## Required Skills

- `spec-driven-development` (state/task management, lane detection)
- `research`
- `agent-browser` (required when implementing frontend work to validate behavior, semantics, and computed CSS)

## Inputs

> [!IMPORTANT]
> Resolve the change set by running `ls changes/ | grep -v archive/`. If exactly one directory exists, use it. Only prompt the user when multiple change sets are present.

## Instructions

1. Load `spec-driven-development` skill and read state from `changes/<name>/state.toml`. Apply state entry check per skill guidelines.

2. Read tasks from `changes/<name>/tasks.toml` if full lane.

3. Determine lane and load plan:
   - **Full lane**: Identify current task (in_progress status). Read corresponding plan from `changes/<name>/plans/`
   - **Vibe/Bug lane**: Read `changes/<name>/plan.md` (single combined plan)

4. **Full lane phase/task entry checks**:
   - Confirm the active phase is `implement`
   - Confirm exactly one task is `in_progress` before execution

5. Execute the plan step by step:
   - Follow steps exactly as written
   - Validate after significant changes
   - Keep repo green
   - Use `research` skill for unexpected code structure or integration questions
   - For any frontend implementation (UI, styling, layout, interaction, accessibility), use `agent-browser` skill to validate the implemented result in a browser
   - During frontend validation, explicitly verify semantic HTML correctness and that computed CSS matches the intended design/reference used during implementation
   - Document any deviations from plan

6. Handle issues:
   - **Minor adjustments**: Proceed and document deviation
   - **Major issues**: Stop and discuss with user
   - **Spec issues** (full lane): Flag for reconciliation

7. Run validation steps from plan, verify acceptance criteria, ensure tests pass.

8. **Full lane completion**:
   - Mark current task complete with `ae sdd task complete [name]` (or `ae sdd task complete --next [name]` when explicitly chaining tasks)
   - Do not manually edit `tasks.toml` task status when CLI task commands are available
   - After task completion, suggest `ae-sdd-next <name>`:
     - If tasks remain, `ae-sdd-next` should deterministically loop `implement -> plan`
     - If tasks are complete, `ae-sdd-next` should advance `implement -> reconcile`
   - Do not update phase status directly in this command

9. **Vibe/Bug lane completion**:
   - Discuss next steps with user
   - If keeping work: do not update phase status; suggest `ae-sdd-next <name>` (optional)
   - If throwing away: done, no state update needed

## Examples

**Full lane implementing a task:**

```text
Input: None (change: "password-reset")
Output: "Loading plan 01.md to implement password validator changes."
       Follows steps: update validator.ts, add reset logic, update tests.
       Validation: All tests pass.
       User: "Looks good."
       Output: "Marked task 1 complete. Three tasks remaining—run ae-sdd-next to loop back to plan for the next task."
```

**Vibe lane quick fix:**

```text
Input: "bug-fix" (user has context with plan.md)
Output: "Following plan.md to patch router/routes.ts."
       Implementation complete, tests pass.
       User: "Great, keep this work."
       Output: "Implementation looks complete. If you want to advance phase, run ae-sdd-next."
```
