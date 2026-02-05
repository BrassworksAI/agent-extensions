---
description: Create implementation tasks from specs (full lane)
---

# Tasks

Create implementation tasks for the change set. Full lane only—vibe/bug lanes skip this and use `/sdd/plan`.

## Required Skills

- `spec-driven-development` (state management, lane detection, task operations)

## Inputs

> [!IMPORTANT]
> Resolve the change set by running `ls changes/ | grep -v archive/`. If exactly one directory exists, use it. Only prompt the user when multiple change sets are present.

## Instructions

1. Load `spec-driven-development` skill. Read state from `changes/<name>/state.toml`. Apply state entry check per skill guidelines.

2. **Lane check**: If lane is `vibe` or `bug`, redirect to `/sdd/plan` per skill guidelines.

3. Read proposal, specs from `changes/<name>/specs/`, and thoughts from `changes/<name>/thoughts/`.

4. This is a dialogue. Before creating tasks, present your breakdown thinking:
   - How you'll map spec requirements to tasks
   - Why you're grouping them this way
   - What task order maintains system stability
   - Ask for feedback on granularity and flow.

5. Create tasks by editing `tasks.toml` to add for each task:
   - Title (short, descriptive)
   - Description (what the task accomplishes)
   - Requirements (mapped from spec lines using EARS syntax)
   - Status: pending

6. **Task ordering principles**:
   - Foundations first (models, types, codegen)
   - Then vertical implementation slices
   - Then integration
   - Then validation

7. **Task constraints**:
   - Every task must be completable in one session
   - Independently testable
   - Leave the system in a committable state

8. When approved, update phase status per skill guidelines and suggest `/sdd/plan <name>`.

## Examples

**Full lane, task breakdown dialogue:**

```text
Input: None (full lane, single change "user-reg")
Output: Present breakdown: "I'll scaffold DB models first, then vertical slices. 3 tasks total."
       User approves. Create tasks in tasks.toml for: foundation, implementation, validation.
```

**Task shape sample:**

```toml
[[task]]
name = "foundation"
title = "Foundation - DB models and migrations"
description = "Add user tables and migrations to support registration."
status = "pending"
requirements = [
  "When a new user is created, the system shall persist profile fields.",
  "When validation fails, the system shall return field errors."
]
```

**Wrong lane detected:**

```text
Input: "vibe-lane" (not full lane)
Output: "Vibe lane should use /sdd/plan instead. Redirecting."
```
