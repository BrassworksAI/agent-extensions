---
description: Draft and refine change proposal (full lane)
argument-hint: <change-name>
---

# Proposal

Draft and refine proposal document for a change set. This command is primarily for **full lane** work.

> **Note**: Vibe lane uses `/sdd:fast:vibe` which creates a lightweight `context.md` instead of a formal proposal. Bug lane uses `/sdd:fast:bug` which handles triage and context creation.

## Arguments

- `$ARGUMENTS` - Change set name

## Instructions

> **SDD Process**: Read `.augment/skills/sdd-state-management.md` for state management guidance.

> **Research**: When needed, delegate to `@librarian` for codebase context. See `.augment/skills/research.md` (project-local) or `~/.augment/skills/research.md` (global) for guidance.

### Setup

1. Read `changes/<name>/state.md`
2. Read `changes/<name>/seed.md` if exists
3. Read `changes/<name>/proposal.md` if exists

### Entry Check

Apply state entry check logic from `.augment/skills/sdd-state-management.md`.

### Research Phase (Recommended)

For non-trivial proposals, delegate to `@librarian`:

1. **Investigate codebase** to understand:
   - Does similar functionality already exist?
   - What would this change interact with?
   - Are there existing patterns or constraints to respect?

2. **Inform proposal** with findings:
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

For vibe or bug work, redirect user to appropriate fast command.

Update state.md `## Lane` with selected lane.

### Collaborative Refinement

This is a dialogue, not a generation task. The user brings whatever definition they have—a half-formed idea, a detailed vision, a conversation summary—and you help them refine it into something complete enough for specs.

**The user is domain expert.** They have tribal knowledge about users, business context, past decisions, and constraints that isn't written anywhere. Your job is to help them externalize and articulate it, not to invent it.

**Ask, don't assume.** When there's a gap, surface it as a question: "What happens if upload fails mid-stream?" Don't fill gaps with guesses.

**Recommend, don't decide.** If you see a pattern in codebase that might apply: "Based on how similar features work, you might want to handle this as X—does that fit?" Let user decide.

**Surface logic gaps.** When described behavior has holes: "You mentioned auto-save, but what happens if user is offline? Is that something we need to handle?" Ask if it should be considered, then let them answer.

Update state.md `## Notes` with key decisions and progress during this phase.

### Proposal Content

A proposal describes **how system works**—written so anyone could read it and understand feature.

**Start with problem.** One or two sentences on what pain exists today. This anchors everything that follows.

**Cover full behavior.** Write in narrative prose. Be specific and concrete.

- What does user do? What happens in response?
- What does system do automatically? When and why?
- What happens when things go wrong? How does system recover or report?
- If there are different modes or states, how do they differ?

**Leave nothing implied.** If there's behavior user is picturing but hasn't written down, it won't make it to specs. The proposal is complete source of truth for what gets built.

**Stay affirmative.** Describe what system does, not what it doesn't do. Scope exclusions, future considerations, and hedged commitments ("if time permits...") dilute focus and create noise for specs phase. If you need to track these, put them in `boundaries.md` or `future.md`—not in proposal.

### Structure Sensing

As user describes feature, notice when natural boundaries emerge:

**Capability boundaries.** "It sounds like there are two distinct capabilities here: edit flow itself, and validation that happens before saving. Am I reading that right?"

**Domain awareness.** "So this would live under notifications, or is it more of a user-preferences thing?" Ask rather than decide—help user identify where this belongs.

Don't force structure prematurely. The goal is to surface what's emerging and check if your framing matches user's mental model. This makes specs phase easier without imposing a taxonomy.

### Readiness

A proposal is ready for specs when you can answer "how does system handle..." for any scenario using only what's written.

### Critique (Recommended)

For full lane proposals, suggest user run `/sdd:tools:critique` for analytical critique:

- Identifies contradictions, gaps, and risks
- Challenges unstated assumptions
- Provides honest assessment of proposal quality

Address any serious issues raised before proceeding.

### Scenario Testing (Recommended)

After critique, suggest user run `/sdd:tools:scenario-test` for user-perspective validation:

- Tests proposal by inhabiting a realistic user persona
- Identifies gaps, friction points, and ambiguities
- Reports whether a user could actually accomplish their goals

Address blocking issues before proceeding; note friction points for consideration.

### Completion

When they explicitly approve proposal:

1. Update state.md: `## Phase Status: complete`, clear `## Notes`
2. Suggest `/sdd:specs <name>` to write change-set specifications (`kind: new|delta`)

Do not log completion in `## Pending` (that section is for unresolved blockers/decisions only).
