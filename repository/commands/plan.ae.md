---
description: Research, plan, and prepare for implementation
---

# Plan

Research the codebase and create an implementation plan. Full lane plans one task at a time; vibe/bug combines discovery + tasking + planning into one pass.

## Required Skills

- `spec-driven-development` (state/task management, lane detection)
- `research`

## Inputs

> [!IMPORTANT]
> Resolve the change set by running `ls changes/ | grep -v archive/`. If exactly one directory exists, use it. Only prompt the user when multiple change sets are present.

## Instructions

1. Load `spec-driven-development` skill and read state from `changes/<name>/state.toml`. Apply state entry check per skill guidelines.

2. Read tasks from `changes/<name>/tasks.toml` if full lane.

3. **Full lane**:
   - Identify current or next task using task status
   - Mark task in_progress in tasks.toml
   - Read task details and all `thoughts/` (discovery outputs)
   - Research codebase for patterns, paths, integration, tests

4. **Vibe/Bug lane**:
   - Read `thoughts/` if they exist
   - Research codebase patterns, changes needed, risks
   - This is discovery + tasking + planning in one pass

5. **Dialogue**: Before writing, summarize research findings. If a path is clear, present your recommendation with reasoning. If you see questions, trade-offs, or risks, discuss them with opinions and recommendations. Collaborate through back-and-forth until direction is right.

6. **Plan contents** (all lanes):
   - Objective/Goal (one sentence, direct)
   - Requirements copied from task/spec as bullets
   - Research findings as bullets: paths, patterns, risks, constraints
   - Steps as an ordered list with exact file paths and concrete changes
   - Validation plan at the end in checkbox format

7. **Validation plan format**:
   - Place at the end of the plan
   - Two checkbox sections:
     - Agent validations (during/post implementation)
     - User validations (optional unless required)

8. **Full lane output**:
   - Create `changes/<name>/plans/<NN>.md`
   - Include all required sections
   - Steps must be exhaustive enough for direct execution

9. **Vibe/Bug lane output**:
   - Create `changes/<name>/plan.md`
   - Include all required sections
   - Include code or pseudocode for complex logic

10. When user explicitly approves, update phase status per skill guidelines and suggest `implement.ae <name>`. Questions/feedback don't count as approval.

## Examples

**Full lane planning with dialogue:**

```text
Input: "password-reset" (full lane, task 2 pending)
Output: "Validated tasks—no in-progress tasks. Started task 2. Read task 2 and thoughts/. 
        Research shows extending auth/login/validator.ts fits DRY. 
        One concern: validator is getting large. Address now or focus on reset?"
       User: "Stay focused on reset. Show me the plan."
       Output: "Created plans/02.md with exhaustive breakdown: file-by-file changes, 
                code snippets for validator.js, auth/routes.js, tests, validation checklist."
```

**Vibe lane with clear recommendation:**

```text
Input: "bug-fix" (vibe lane)
Output: "Read thoughts/. Found consistent fix pattern in router/middleware.ts. 
        Applying that keeps codebase uniform. Created plan.md with exhaustive breakdown: 
        exact changes to middleware.ts, code snippet, test updates, validation steps."
```
