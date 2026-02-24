---
description: Create implementation tasks from specs (full lane)
---

# Tasks

Create implementation tasks for the change set. Full lane only—vibe/bug lanes skip this and use `plan.ae`.

## Required Skills

- `spec-driven-development` (state management, lane detection, task operations)

## Inputs

> [!IMPORTANT]
> Resolve the change set by running `ls changes/ | grep -v archive/`. If exactly one directory exists, use it. Only prompt the user when multiple change sets are present.

## Instructions

1. Load `spec-driven-development` skill. Read state from `changes/<name>/state.toml`. Apply state entry check per skill guidelines.

2. **Lane check**: If lane is `vibe` or `bug`, redirect to `plan.ae` per skill guidelines.

3. Read proposal, specs from `changes/<name>/specs/`, and thoughts from `changes/<name>/thoughts/`.

4. This is a dialogue. Before creating tasks, present your breakdown thinking:
   - How you'll map spec requirements to tasks
   - Why you're grouping them this way
   - What task order maintains system stability
   - Ask for feedback on granularity and flow.

5. Create tasks by editing `tasks.toml` to add for each task:
   - Title (short, descriptive)
   - Description (what the task accomplishes)
   - `spec_requirements` entries grouped by source spec file
   - Requirement strings mapped from spec lines using EARS syntax
   - Status: pending

   Required structure:

   ```toml
   [[task]]
   name = "task-name"
   title = "Task title"
   description = "What this task accomplishes"
   status = "pending"

   [[task.spec_requirements]]
   spec = "changes/<name>/specs/<file>.md"
   requirements = [
     "When <condition>, the system SHALL <behavior>",
     "If <trigger>, the system SHALL <behavior>"
   ]
   ```

6. **Task ordering principles**:
   - Foundations first (models, types, codegen)
   - Then vertical implementation slices
   - Then integration
   - Then validation

7. **Task constraints**:
   - Every task must be completable in one session
   - Independently testable
   - Leave the system in a committable state

8. Do not update phase status in this command. When the user wants to proceed, suggest `next.ae <name>`.

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

[[task.spec_requirements]]
spec = "changes/user-reg/specs/registration.md"
requirements = [
  "When a new user is created, the system SHALL persist profile fields.",
  "When validation fails, the system SHALL return field errors."
]

[[task.spec_requirements]]
spec = "changes/user-reg/specs/security.md"
requirements = [
  "If password policy fails, the system SHALL reject account creation."
]
```

**Wrong lane detected:**

```text
Input: "vibe-lane" (not full lane)
Output: "Vibe lane should use plan.ae instead. Redirecting."
```
