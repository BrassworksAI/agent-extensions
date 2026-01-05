---
description: Draft and refine change proposal (full lane)
argument-hint: <change-name>
---

# Proposal

Draft and refine the proposal document for a change set. This command is primarily for **full lane** work.

> **SDD Process**: Read `changes/<name>/state.md` first to verify phase is `proposal` (or `ideation` transitioning). If unsure how state management works, read `augment/skills/sdd-state-management.md`.

> **Note**: Vibe lane uses `/sdd:fast:vibe` which creates a lightweight `context.md` instead of a formal proposal. Bug lane uses `/sdd:fast:bug` which handles triage and context creation.

## Arguments

- `$ARGUMENTS` - Change set name

## Instructions

### Setup

1. Read `changes/<name>/state.md` - verify phase is `proposal` (or `ideation` transitioning)
2. Read `changes/<name>/seed.md` if exists (for context)
3. Read `changes/<name>/proposal.md` if exists

### Research Phase (Recommended)

For non-trivial proposals, delegate to `@librarian`:

1. **Investigate the codebase** to understand:
   - Does similar functionality already exist?
   - What would this change interact with?
   - Are there existing patterns or constraints to respect?

2. **Inform the proposal** with findings:
   - Reference existing code/patterns in approach
   - Note integration points
   - Identify potential risks based on codebase structure

### Lane Selection

If lane not yet selected, determine with user:

| Lane | When to Use |
|------|-------------|
| `full` | New capabilities, architectural changes, complex features |
| `vibe` | Prototypes, experiments, quick enhancements (use `/sdd:fast:vibe` instead) |
| `bug` | Fixing incorrect behavior (use `/sdd:fast:bug` instead) |

For vibe or bug work, redirect user to the appropriate fast command.

Update state.md with selected lane.

### Proposal Content

Proposals are **freeform** - capture intent in whatever structure works. Common elements:

- **Problem**: What problem are we solving?
- **Goals**: What does success look like?
- **Non-Goals**: What are we explicitly NOT doing?
- **Approach**: High-level solution direction
- **Risks**: What could go wrong?
- **Definition of Done**: How do we know we're finished?

### Critique (Recommended)

For full lane proposals, suggest the user run `/sdd:tools:critique` for analytical critique:

- Identifies contradictions, gaps, and risks
- Challenges unstated assumptions
- Provides honest assessment of proposal quality

Address any serious issues raised before proceeding.

### Scenario Testing (Recommended)

After critique, suggest the user run `/sdd:tools:scenario-test` for user-perspective validation:

- Tests proposal by inhabiting a realistic user persona
- Identifies gaps, friction points, and ambiguities
- Reports whether a user could actually accomplish their goals

Address blocking issues before proceeding; note friction points for consideration.

### Completion

Work through the proposal collaboratively with the user. When they explicitly approve:

1. Log approval in state.md under `## Pending`:
   ```
   None - Proposal approved: [brief summary of agreed approach]
   ```
2. Update state.md phase to `specs`
3. Suggest `/sdd:specs <name>` to write delta specifications

Don't advance until the user clearly signals approval. Questions, feedback, or acknowledgments don't count as approval.
