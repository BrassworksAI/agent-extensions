---
description: Create implementation tasks from specs (full lane)
---

# Tasks

Create implementation tasks for the change set. Full lane only—vibe/bug lanes skip this and use `/sdd/plan`.

## Inputs

> [!IMPORTANT]
> Resolve the change set by running `ls changes/ | grep -v archive/`. If exactly one directory exists, use it. Only prompt the user when multiple change sets are present.

## Instructions

Load `spec-driven-development` skill. Read state, proposal, specs, and thoughts from the change set. If lane is `vibe` or `bug`, redirect to `/sdd/plan`.

This is a dialogue. Before writing tasks.md, present your breakdown thinking: how you'll map spec requirements, why you're grouping them this way, and what task order maintains system stability. Ask for feedback on granularity and flow. Update state.md `## Notes` with decisions.

Create `changes/<name>/tasks.md` with checkbox tracking: `[ ]` pending, `[o]` in progress, `[x]` complete. Do not number tasks. Each task includes a short title, a description, and requirements mapped to spec lines. Order tasks: foundations first (models, types, codegen), vertical implementation slices, integration, validation. Every task must be completable in one session, independently testable, and leave the system in a committable state.

## Examples

**Full lane, single tasks.md:**

```text
Input: None (full lane, single change "user-reg")
Output: Present breakdown: "I'll scaffold DB models first, then vertical slices. 3 tasks total."
       User approves. Create tasks.md with [ ] foundation, [ ] implementation, [ ] validation.
```

**Task shape sample (no numbering):**

```text
[ ] Title: Foundation - DB models and migrations
    Description: Add user tables and migrations to support registration.
    Requirements:
    - When a new user is created, the system shall persist profile fields.
    - When validation fails, the system shall return field errors.
```

**Wrong lane detected:**

```text
Input: "vibe-lane" (not full lane)
Output: "Vibe lane should use /sdd/plan instead. Redirecting."
```
