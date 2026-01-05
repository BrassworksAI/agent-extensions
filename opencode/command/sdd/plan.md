---
name: sdd/plan
description: Research, plan, and prepare for implementation
agent: sdd/forge
---

<skill>sdd-state-management</skill>
<skill>research</skill>

# Plan

Research the codebase and create an implementation plan.

## Arguments

- `$ARGUMENTS` - Change set name (optionally: `<name> <task-number>` for full lane)

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
2. Find first pending task (or specified task number)
3. Read any existing plans in `changes/<name>/plans/`

### Research

Use the `research` skill to understand:
- Where changes need to happen (exact file paths)
- Existing patterns to follow
- Integration points
- Tests that need updates

### Create Plan

Create `changes/<name>/plans/<NN>.md`:

```markdown
# Plan: Task <N> - <Task Title>

## Objective

<From task description>

## Requirements

- "<EARS requirement line from delta specs>"

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

- [ ] Acceptance criteria
- [ ] Tests pass
```

### Completion

Work through the plan collaboratively with the user. When they explicitly approve:

1. Mark task as `in_progress` in tasks.md
2. Log approval in state.md under `## Pending`:
   ```
   None - Plan approved for task [N]: [task title]
   ```
3. Update state.md phase to `implement`
4. Suggest `/sdd/implement <name>`

Don't advance until the user clearly signals approval. Questions, feedback, or acknowledgments don't count as approval.

---

## Vibe/Bug Lane Planning

For vibe/bug lanes, combine discovery + tasking + planning into one pass. Get to building fast.

### Research

Use the `research` skill to understand:
- What exists in the codebase
- Where changes need to happen
- Patterns to follow
- Potential risks or complications

### Create Combined Plan

Create `changes/<name>/plan.md` (single file, not per-task):

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

Work through the plan collaboratively with the user. When they explicitly approve:

1. Log approval in state.md under `## Pending`:
   ```
   None - Plan approved: [brief summary of approach]
   ```
2. Update state.md phase to `implement`
3. Suggest `/sdd/implement <name>`

Don't advance until the user clearly signals approval. Questions, feedback, or acknowledgments don't count as approval.

---

## Plan Quality

A good plan:
- Is grounded in research (not assumptions)
- Has specific file paths (verified to exist)
- Follows patterns found in the codebase
- Includes validation steps
- Is appropriately detailed for the lane
