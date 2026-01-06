---
description: Research, plan, and prepare for implementation
argument-hint: <change-name> [task-number]
---

# Plan

Research the codebase and create an implementation plan.

> **SDD Process**: Read `changes/<name>/state.md` first to verify phase is `plan`. If unsure how state management works, read `augment/skills/sdd-state-management.md`.

> **Research**: When you need to investigate the codebase, delegate to `@librarian` with a specific research question. See `augment/skills/research.md` for guidance.

## Arguments

- `$ARGUMENTS` - Change set name (for full lane, plans the next incomplete task)

## Instructions

### Setup

1. Read `changes/<name>/state.md` - verify phase is `plan`
2. Determine lane (full, vibe, or bug)
3. Read context:
   - **Full lane**: Read `tasks.md`, identify current task
   - **Vibe/Bug lane**: Read `context.md`

---

## Full Lane Planning

For full lane, plan one task at a time.

### Identify Current Task

1. Read `changes/<name>/tasks.md`
2. If a task is already `[o]` (In Progress), assume it is the current task and continue planning it.
3. Otherwise, pick the first `[ ]` (Pending) task as the current task.
4. Read any existing plans in `changes/<name>/plans/`

### Load Specs (Full Lane)

Read all change-set specs for context, recursively:
- `changes/<name>/specs/**/*.md`

Do not assume specs only exist at `changes/<name>/specs/*.md`.

### Research

Delegate to `@librarian` to understand:
- Where changes need to happen (exact file paths)
- Existing patterns to follow
- Integration points
- Tests that need updates

If you might introduce new dependencies, APIs, or framework features, also do quick online research to confirm current best practices and version-specific details.

### Collaborative Planning

This command is a **dialogue**, not a one-way generation.

1. **Think Out Loud**: Before writing any plan file, summarize what you learned from research and the approach you’re leaning toward.
2. **Present Options**: If there are multiple reasonable approaches, present 2–3 with trade-offs.
3. **Ask Targeted Questions**: Clarify any unknowns that materially affect the plan (scope boundaries, rollout strategy, time/complexity constraints).
4. **Invite Feedback**: Ask the user if they want any changes to structure, granularity, or sequencing.
5. **Iterate**: Only write the plan file once the user indicates the direction looks right.

### Create Plan

At the start of planning (before writing the plan file), update `changes/<name>/tasks.md` to mark the current task as `[o]` (In Progress). Ensure there is at most one `[o]` task at a time.

Only after the user indicates alignment, create `changes/<name>/plans/<NN>.md` (or update the existing plan for the in-progress task):

```markdown
# Plan: <Task Title>

## Objective

<From task description>

## Requirements

- **From <spec>**: "<full EARS requirement line from change-set specs>"

## Research Findings

- <Patterns found>
- <Integration points>
- <Files to modify>

## Steps

### Step 1: <Title>

**Files:** `path/to/file.ts`

**Changes:**
- Specific modifications

### Step 2: <Title>

...

## Validation

- [ ] Acceptance criteria (copy from tasks.md)
- [ ] Tests pass
```

### Completion

After the plan file is written, review it with the user. When they explicitly approve:

1. Update `changes/<name>/state.md` phase to `implement`
2. If an approval record is needed, capture it in a separate artifact (e.g., `changes/<name>/thoughts/decisions.md`), not `## Pending`.
3. Suggest `/sdd:implement <name>`

Don't advance until the user clearly signals approval. Questions, feedback, or acknowledgments don't count as approval.

---

## Vibe/Bug Lane Planning

For vibe/bug lanes, combine discovery + tasking + planning into one pass. Get to building fast.

### Research

Delegate to `@librarian` to understand:
- What exists in the codebase
- Where changes need to happen
- Patterns to follow
- Potential risks or complications

### Collaborative Planning

This command is a **dialogue**, not a one-way generation.

1. **Think Out Loud**: Before writing any plan file, summarize what you learned from research and the approach you’re leaning toward.
2. **Present Options**: If there are multiple reasonable approaches, present 2–3 with trade-offs.
3. **Ask Targeted Questions**: Clarify any unknowns that materially affect the plan.
4. **Invite Feedback**: Ask the user if they want to adjust scope or sequencing.
5. **Iterate**: Only write the plan file once the user indicates the direction looks right.

### Create Combined Plan

Only after the user indicates alignment, create `changes/<name>/plan.md` (single file, not per-task):

```markdown
# Plan: <name>

## Goal

<What we're trying to accomplish - from context.md>

## Research Findings

<What we learned about the codebase>
- Relevant files and patterns
- Integration points
- Potential risks

## Approach

<High-level strategy>

## Changes

### 1. <First change>

**Files:** `path/to/file.ts`

**What:**
- Specific modifications

### 2. <Second change>

...

## Validation

How to verify it works:
- [ ] Test case 1
- [ ] Test case 2
```

### Keep It Lean

For vibe/bug lanes:
- Don't over-plan - you're exploring
- Enough detail to start, not a complete spec
- If it gets complicated, that's a sign to consider full lane

### Completion

After the plan file is written, review it with the user. When they explicitly approve:

1. Update `changes/<name>/state.md` phase to `implement`
2. If an approval record is needed, capture it in a separate artifact (e.g., `changes/<name>/thoughts/decisions.md`), not `## Pending`.
3. Suggest `/sdd:implement <name>`

Don't advance until the user clearly signals approval. Questions, feedback, or acknowledgments don't count as approval.

---

## Plan Quality

A good plan:
- Is grounded in research (not assumptions)
- Has specific file paths (verified to exist)
- Follows patterns found in the codebase
- Includes validation steps
- Is appropriately detailed for the lane
