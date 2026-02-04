---
description: Research, plan, and prepare for implementation
---

# Plan

Research the codebase and create an implementation plan. Full lane plans one task at a time; vibe/bug combines discovery + tasking + planning into one pass.

## Inputs

> [!IMPORTANT]
> Resolve the change set by running `ls changes/ | grep -v archive/`. If exactly one directory exists, use it. Only prompt the user when multiple change sets are present.

## Instructions

Load `spec-driven-development`, `research`, and `keep-current` skills. Read state.md and tasks.md. Apply state entry check.

**Full lane**: Identify current task. If `[o]` exists, plan it (never replace without user approval). Otherwise, pick first `[ ]`, mark it `[o]`. Read task from tasks.md (which contains the spec) and all `thoughts/` (discovery outputs). Research codebase for patterns, paths, integration, tests. Update state.md notes.

**Vibe/bug lane**: Read `thoughts/` and research codebase patterns, changes needed, risks.

This is a dialogue. Before writing, summarize research findings. If a path is clear, present your recommendation with reasoning. If you see questions, trade-offs, or risks, discuss them with opinions and recommendations. Collaborate through back-and-forth until direction is right.

**Plan contents (all lanes)**:

- Objective/Goal (one sentence, direct).
- Requirements copied from task/spec as bullets.
- Research findings as bullets: paths, patterns, risks, and constraints.
- Steps as an ordered list with exact file paths and concrete changes.
- Validation plan at the end in checkbox format.

**Validation plan format (all lanes)**:

- Place at the end of the plan.
- Two checkbox sections:
  - Agent validations (during/post install).
  - User validations (optional unless the agent deems user-run validation required).

**Full lane output**:

- Create `changes/<name>/plans/<NN>.md`.
- Include all required sections listed above.
- Steps must be exhaustive enough for direct execution without interpretation.

**Vibe/bug lane output**:

- Create `changes/<name>/plan.md`.
- Include all required sections listed above.
- Include code or pseudocode for complex logic.

When user explicitly approves, update state.md: `## Phase Status: complete`, clear `## Notes`, suggest `/sdd/implement <name>`. Questions/feedback don't count as approval.

## Examples

**Full lane planning with dialogue:**

```text
Input: "password-reset" (full lane, task 2 pending)
Output: "Validated tasks.md—no in-progress tasks. Marked task 2 in-progress. Read task 2 and thoughts/. Research shows extending auth/login/validator.ts fits DRY. One concern: validator is getting large. Address now or focus on reset?"
       User: "Stay focused on reset. Show me the plan."
       Output: "Created plans/02.md with exhaustive breakdown: file-by-file changes, code snippets for validator.js, auth/routes.js, tests, validation checklist."
```

**Vibe lane with clear recommendation:**

```text
Input: "bug-fix" (vibe lane)
Output: "Read thoughts/. Found consistent fix pattern in router/middleware.ts. Applying that keeps codebase uniform. Created plan.md with exhaustive breakdown: exact changes to middleware.ts, code snippet, test updates, validation steps."
```
