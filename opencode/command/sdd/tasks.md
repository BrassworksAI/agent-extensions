---
name: sdd/tasks
description: Create implementation tasks from specs (full lane)
agent: sdd/forge
---

<skill>sdd-state-management</skill>
<skill>spec-format</skill>

# Tasks

Create implementation tasks for the change set. This command is for **full lane** only.

> **Note**: Vibe and bug lanes skip this command. They use `/sdd/plan` which combines research, tasking, and planning into a single `plan.md` file.

## Arguments

- `$ARGUMENTS` - Change set name

## Instructions

### Setup

1. Read `changes/$1/state.md` - verify phase is `tasks` and lane is `full`

If lane is `vibe` or `bug`, redirect user to `/sdd/plan` instead.

## Context

### Proposal
!`cat changes/$1/proposal.md 2>/dev/null || echo "No proposal found"`

### Specs
!`find changes/$1/specs -name "*.md" -exec cat {} + 2>/dev/null || echo "No specs found"`

### Thoughts
!`find changes/$1/thoughts -name "*.md" -exec cat {} + 2>/dev/null || echo "No thoughts found"`

## Instructions

### Collaborative Tasking

This command is a **dialogue**, not a one-way generation.

1. **Think Out Loud**: Before writing the file, present your initial thoughts on the task breakdown.
   - Generalize the requirements: Tasks should include *anything* necessary to make the implementation successful (e.g., scaffolding, environment setup, refactoring).
   - Traceability: You MUST explicitly call out which specific spec requirements are covered by each task.
   - Explain *why* you're grouping certain requirements and why you've chosen a specific order (e.g., maintaining system stability, foundation-first).
2. **Present Options**: If there are multiple valid ways to slice the work (e.g., horizontal vs vertical, foundation-first vs feature-first), present them to the user with trade-offs.
3. **Invite Feedback**: Explicitly ask the user if they have specific preferences for task granularity or if there's a specific logical flow they want to follow to maintain system stability.
4. **Iterate**: Only write `tasks.md` once a consensus on the strategy has been reached.

### Task Structure

Create `changes/<name>/tasks.md` using checkbox-style progress tracking:
- `[ ]` = Pending
- `[o]` = In Progress (exactly one task at a time)
- `[x]` = Complete

```markdown
# Tasks: <name>

## Overview

Brief summary of what these tasks accomplish.

## Tasks

### [ ] <Title>

**Description:**
What this task accomplishes. Focus on why it exists and what it changes.

**Requirements:**

#### Foundations & Prerequisites (if any, else skip)
- <Description of technical prerequisite, scaffolding, or environment setup needed for success>

#### <spec-path>
- "<full EARS requirement line>"
...
```

### Task Ordering & Logic

Design tasks to ensure the application is **never in a broken state** and can be **committed after every task**.

1. **Foundations First**: Models, types, interfaces, and database migrations.
2. **Implementation Slices**: Implement functionality in vertical slices that introduce new code paths behind flags or as new modules before wiring them in.
3. **Integration**: Connect new components to existing systems.
4. **Validation**: Test suites, cleanup of old paths, and consolidation.

Order tasks by dependency. A task is only "done" when the system is stable and committable.

### Task Granularity

Each task should be:
- Completable in one implementation session
- Independently testable
- Clear on what "done" means

### Requirement Mapping

- Every requirement in delta specs must map to at least one task
- Tasks MUST reference requirements by quoting the full EARS line and specifying the source spec file
- Use `spec-format` skill to understand requirement structure

### Completion

Work through the task breakdown collaboratively with the user. 

1. **Be Proactive**: If the grouping or ordering is ambiguous, propose a logical structure based on architectural best practices ("Foundations First", "Implementation Slices"). Explain the rationale behind your choices.
2. **Ask for Clarification**: If a requirement is genuinely ambiguous or could be grouped in multiple valid ways with different trade-offs, ask the user for their preference.
3. **Approval**: When the user explicitly approves:
   - Log approval in state.md under `## Pending`:
     ```
     None - Tasks approved: [number] tasks defined
     ```
   - Update state.md phase to `plan`
   - Suggest running `/sdd/plan <name>` to plan first task

Don't advance until the user clearly signals approval. Questions, feedback, or acknowledgments don't count as approval.
