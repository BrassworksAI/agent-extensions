---
name: sdd/brainstorm
description: Explore problem space and develop seed document
agent: sdd/plan
---

# Required Skills (Must Load)

You MUST load and follow these skills before doing anything else:

- `sdd-state-management`
- `research`

If any required skill content is missing or not available in context, you MUST stop and ask the user to re-run the command or otherwise provide the missing skill content. Do NOT proceed without them.

# Brainstorm / Ideation

Explore the problem space collaboratively to develop a seed document.

## Arguments

- `$ARGUMENTS` - Change set name

## Instructions

### Setup

!`cat changes/$1/state.md 2>/dev/null || echo "State file not found"`

!`cat changes/$1/seed.md 2>/dev/null || echo "No seed found"`

### Entry Check

Apply state entry check logic from `sdd-state-management` skill.

### Research Phase (As Needed)

During ideation, research can help ground ideas in reality:

1. **Use `research` skill** when you need to understand:
   - Does something similar already exist in the codebase?
   - What constraints does the current architecture impose?
   - What patterns are already established?

2. **Use research to inform ideation**, not constrain it:
   - Research helps identify what's possible
   - But don't let existing patterns limit creative thinking
   - Note constraints, then explore solutions

### Ideation Process

This is a **collaborative conversation**. Your job is to:

1. **Understand the problem**: Ask clarifying questions about what the user wants to build
2. **Explore constraints**: What are the boundaries? What's out of scope?
3. **Surface assumptions**: What's being taken for granted?
4. **Identify risks**: What could go wrong?
5. **Document incrementally**: Update seed.md and state.md `## Notes` as understanding develops

### Seed Document Structure

The seed is freeform but typically includes:

```markdown
# <Name> Seed

## Problem Statement

What problem exists and why it matters.

## Core Thesis

Key beliefs/assumptions underlying the solution.

## Proposed Approach

High-level direction (not detailed design).

## Constraints

What limits the solution space.

## Open Questions

What needs to be resolved.

## Risks

What could derail this.
```

### Critique

When the seed feels complete, suggest the user run `/sdd/tools/critique seed`:

- Identifies contradictions and missing cases
- Flags risks that aren't acknowledged
- Validates the seed is ready to become a proposal

If critique identifies serious issues, work with user to address them.

### Completion

When seed is solid and user explicitly approves:

1. Update state.md: `## Phase Status: complete`, clear `## Notes`
2. Suggest running `/sdd/proposal <name>`

Do not log completion in `## Pending` (that section is for unresolved blockers/decisions only).
