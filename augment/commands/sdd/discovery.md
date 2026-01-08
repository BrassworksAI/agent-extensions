---
description: Discover high-level architectural requirements for change-set specs
argument-hint: <change-set-name>
---

# Discovery

Understand high-level architectural requirements for implementing change-set specs. This phase answers big questions about how change fits into—or extends—the existing architecture.

## Purpose

Discovery is NOT about planning implementation details. It's about:
- Understanding what architectural patterns/systems specs will touch
- Identifying if change slots cleanly into existing architecture (simple case)
- Recognizing when architectural concerns need resolution before planning (complex case)
- Working through high-level solutions when path isn't obvious

## Arguments

- `$ARGUMENTS` - Change set name

## Instructions

> **SDD Process**: Read `.augment/skills/sdd-state-management.md` for state management guidance.

> **Research**: When needed, delegate to `@librarian` for codebase context. See `.augment/skills/research.md` (project-local) or `~/.augment/skills/research.md` (global) for guidance.

> **Architecture**: Use guidance from `.augment/skills/architecture-fit-check.md` and `.augment/skills/architecture-workshop.md` (project-local) or `~/.augment/skills/architecture-fit-check.md` and `~/.augment/skills/architecture-workshop.md` (global).

### Setup

1. Read `changes/<name>/state.md`
2. Read all change-set specs from `changes/<name>/specs/**/*.md` (both `kind: new` and `kind: delta`)
3. Read `changes/<name>/proposal.md` for context

### Entry Check

Apply state entry check logic from `.augment/skills/sdd-state-management.md`.

If lane is not `full`, redirect user to appropriate command.

### Research Phase (Critical)

Before evaluating architecture fit, delegate to `@librarian` to understand codebase:

1. **Research to understand**:
   - Current architecture patterns in codebase
   - Code areas that will be affected by specs
   - Existing implementations of similar capabilities
   - Potential integration points and conflicts

2. **Build context** for architecture evaluation:
   - Document what you learned about architecture
   - Identify specific code areas specs will touch
   - Note any patterns that seem relevant

Update state.md `## Notes` with architectural findings and research insights.

### Architecture Assessment

Using architecture-fit-check framework, answer primary question:

**Can these change-set specs be implemented cleanly within existing architecture?**

#### Simple Case: Clean Fit

If specs slot easily into existing architecture (e.g., new endpoint, data to template, straightforward CRUD), there's not much to record here:
- Note that architecture review found no concerns
- Proceed directly to tasks phase

#### Complex Case: Concerns Exist

If specs WOULD work but raise concerns:
- Would require messy workarounds
- Introduces inconsistent patterns
- Creates technical debt
- Requires primitives codebase doesn't have

Then adopt **Daedalus personality** (master architect) to work through best solution with user.

### Daedalus Mode (When Concerns Exist)

When architectural concerns are identified, engage as Daedalus—the master architect who designs elegant solutions:

1. **Explain concern clearly** to user:
   - What makes straightforward approach problematic
   - Why it matters for maintainability/consistency
   - What questions need answering

2. **Explore high-level solutions**:
   - **Light-touch options**: Adapter layer, new module boundary, small abstraction
   - **Architecture options**: New eventing/pubsub system, state management pattern, concurrency model

3. **Work through it with user**:
   - Present tradeoffs (blast radius, incremental path, long-term impact)
   - Get user input on direction
   - Reach consensus on approach

4. **Capture thoughts along way** in `changes/<name>/thoughts/`:
   - Create files as needed during session
   - Free-form format—whatever captures exploration
   - Document concerns, options considered, decisions reached
   - This preserves context if user continues in a new chat

### Thoughts Directory

Discovery outputs to `changes/<name>/thoughts/`. This is a free-form workspace:

```
changes/<name>/
  thoughts/
    architecture-concerns.md
    options-explored.md
    decision-rationale.md
```

Create as many files as needed. The goal is capturing architectural exploration so it's not lost.

Update state.md `## Notes` with architecture decisions and rationale.

### Updating Specs

If discovery reveals specs themselves need changes:
- Return to specs phase
- Update change-set specs to reflect architectural decisions
- Re-run discovery

### Completion

When they explicitly approve architecture findings:

1. Update state.md: `## Phase Status: complete`, clear `## Notes`
2. Suggest running `/sdd:tasks <name>`

Do not log completion in `## Pending` (that section is for unresolved blockers/decisions only).
