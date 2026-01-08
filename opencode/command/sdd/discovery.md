---
name: sdd/discovery
description: Discover high-level architectural requirements for change-set specs
agent: sdd/plan
---

<skill>sdd-state-management</skill>
<skill>research</skill>
<skill>architecture-fit-check</skill>
<skill>architecture-workshop</skill>

# Discovery

Understand the high-level architectural requirements for implementing change-set specs. This phase answers the big questions about how the change fits into—or extends—the existing architecture.

## Purpose

Discovery is NOT about planning implementation details. It's about:
- Understanding what architectural patterns/systems the specs will touch
- Identifying if the change slots cleanly into existing architecture (simple case)
- Recognizing when architectural concerns need resolution before planning (complex case)
- Working through high-level solutions when the path isn't obvious

## Arguments

- `$ARGUMENTS` - Change set name

## Instructions

### Setup

!`cat changes/$1/state.md 2>/dev/null || echo "State file not found"`

!`cat changes/$1/proposal.md 2>/dev/null || echo "No proposal found"`

!`find changes/$1/specs -name "*.md" -exec cat {} + 2>/dev/null || echo "No specs found"`

### Entry Check

Apply state entry check logic from `sdd-state-management` skill.

If lane is not `full`, redirect user to appropriate command.

### Research Phase (Critical)

Before evaluating architecture fit, use the `research` skill to understand the codebase:

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

Using `architecture-fit-check` skill framework, answer the primary question:

**Can these change-set specs be implemented cleanly within the existing architecture?**

#### Simple Case: Clean Fit

If the specs slot easily into existing architecture (e.g., new endpoint, data to template, straightforward CRUD), there's not much to record here:
- Note that architecture review found no concerns
- Proceed directly to tasks phase

#### Complex Case: Concerns Exist

If the specs WOULD work but raise concerns:
- Would require messy workarounds
- Introduces inconsistent patterns
- Creates technical debt
- Requires primitives the codebase doesn't have

Then adopt the **Daedalus personality** (master architect) to work through the best solution with the user.

### Daedalus Mode (When Concerns Exist)

When architectural concerns are identified, engage as Daedalus—the master architect who designs elegant solutions:

1. **Explain the concern clearly** to the user:
   - What makes the straightforward approach problematic
   - Why it matters for maintainability/consistency
   - What questions need answering

2. **Explore high-level solutions**:
   - **Light-touch options**: Adapter layer, new module boundary, small abstraction
   - **Architecture options**: New eventing/pubsub system, state management pattern, concurrency model

3. **Work through it with the user**:
   - Present tradeoffs (blast radius, incremental path, long-term impact)
   - Get user input on direction
   - Reach consensus on approach

4. **Capture thoughts along the way** in `changes/<name>/thoughts/`:
   - Create files as needed during the session
   - Free-form format—whatever captures the exploration
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

Create as many files as needed. The goal is capturing the architectural exploration so it's not lost.

Update state.md `## Notes` with architecture decisions and rationale.

### Updating Specs

If discovery reveals the specs themselves need changes:
- Return to specs phase
- Update change-set specs to reflect architectural decisions
- Re-run discovery

### Completion

When they explicitly approve the architecture findings:

1. Update state.md: `## Phase Status: complete`, clear `## Notes`
2. Suggest running `/sdd/tasks <name>`

Do not log completion in `## Pending` (that section is for unresolved blockers/decisions only).
