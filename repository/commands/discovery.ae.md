---
description: Discover high-level architectural requirements for change-set specs
---

# SDD Discovery

Analyze high-level architectural requirements for implementing change-set specs to ensure proper alignment and identify potential concerns.

## Required Skills

- `spec-driven-development`
- `research`
- `architecture-fit-check`
- `architecture-workshop`

## Inputs

> [!IMPORTANT]
> Resolve the change set by running `ls changes/ | grep -v archive/`. If exactly one directory exists, use it. Only prompt the user when multiple change sets are present.

## Instructions

Discovery answers how a change fits into or extends the existing architecture.

1. **Setup**: Run `cat changes/<name>/state.md`, `proposal.md`, and any existing specs in `changes/<name>/specs`.
2. **Entry Check**: Apply state entry check logic from `spec-driven-development`. If the lane is not `full`, redirect the user.
3. **Research Phase**: Use the `research` skill to understand current architectural patterns, affected code areas, and existing implementations. Update `state.md` with findings.
4. **Architecture Assessment**: Use `architecture-fit-check`. If it's a "Clean Fit," note it and move to the tasks phase.
5. **Daedalus Mode (Complex Case)**: If concerns exist (e.g., technical debt, messy workarounds), adopt the master architect persona. Explain concerns, explore solutions (light-touch vs. architectural), and reach consensus with the user.
6. **Capture Thoughts**: Document explorations, tradeoffs, and decisions in `changes/<name>/thoughts/`.

## Success Criteria

- Architectural findings and research insights are documented in `state.md` notes.
- If complex, decision rationale is captured in `changes/<name>/thoughts/`.
- `state.md` is updated to `## Phase Status: complete`.

## Usage Examples

### Do

- "This change slots cleanly into the existing `Notification` pattern."
- "Implementing this directly would create circular dependencies. Should we introduce an event bus instead?"

### Don't

- Start planning implementation tasks or file-level changes.
- Guess at architectural alignment without researching the codebase.
