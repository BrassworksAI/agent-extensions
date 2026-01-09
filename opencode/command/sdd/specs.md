---
name: sdd/specs
description: Write change-set specifications for change
agent: sdd/plan
---

# Required Skills (Must Load)

You MUST load and follow these skills before doing anything else:

- `sdd-state-management`
- `spec-format`
- `research`

If any required skill content is missing or not available in context, you MUST stop and ask the user to re-run the command or otherwise provide the missing skill content. Do NOT proceed without them.

# Specs

Write change-set specifications for the change set (`kind: new` and `kind: delta`).

## Arguments

- `$ARGUMENTS` - Change set name

## Instructions

### Setup (Injected Context)

The following commands print the authoritative documents into the chat as context:

!`cat changes/$1/state.md 2>/dev/null || echo "State file not found"`

!`cat changes/$1/proposal.md 2>/dev/null || echo "No proposal found"`

You MUST treat the printed contents as the primary source of truth for this run.

- Do NOT waste cycles re-searching for `changes/$1/state.md` or `changes/$1/proposal.md` when they are already shown.
- Do NOT ask the user to paste them again unless the output says they were not found.
- You MUST update `changes/$1/state.md` as instructed by `sdd-state-management` (notes, phase status, pending decisions).

### Entry Check

Apply state entry check logic from `sdd-state-management` skill.

If lane is not `full`, redirect user to appropriate command.

### Collaborative Discovery (Required)

Spec writing is a **collaborative** process. You MUST NOT immediately start writing specification files.

Before any spec text is written, you MUST run a brief **collaborative alignment loop**:

1. **Domain understanding (show your work)**
   - Summarize the user’s goal, actors, workflows, and constraints in your own words.
   - Identify assumptions and explicitly mark them as assumptions.

2. **Capability taxonomy understanding (map intent to structure)**
   - Explain how you believe the change maps into the existing capability hierarchy.
   - Call out any unclear boundaries (what is in-scope vs out-of-scope).

3. **Back-and-forth confirmation (required)**
   - Present your understanding as a set of concrete statements and decisions.
   - Offer a small set of options for any open decision (with tradeoffs).
   - Pause and wait for the user to confirm or correct.

Only after the user explicitly confirms (or corrects) the understanding above may you proceed to Research/Taxonomy/Spec writing.

### Research Phase

After the collaborative discovery is confirmed, use the `research` skill:

1. **Research to understand**:
   - Current spec structure and taxonomy
   - Related existing capabilities
   - How similar things are specified

2. **Build context** for spec writing:
   - What specs already exist in related areas
   - What naming conventions are used

### Taxonomy Mapping

With research in hand, suggest the user run `/sdd/tools/taxonomy-map <name>` and confirm the result with them:

- Determines where new capabilities should live in the spec hierarchy
- Recommends brownfield (existing specs) vs greenfield (new specs)
- Provides boundary decisions and group structure

### Writing Change Set Specs

Only after the user explicitly approves the domain summary and taxonomy mapping:

Create specs in `changes/<name>/specs/` following the `spec-format` skill. Specs may be nested by domain/subdomain under that folder (e.g. `changes/<name>/specs/auth/login.md`).
Remember that change set specs have YAML frontmatter `kind: new | delta`.

1. **Identify capabilities** needed from the proposal
2. **Determine paths** using taxonomy mapping guidance
3. **Write requirements** using EARS syntax

Update state.md `## Notes` with progress on spec writing, key decisions, and any issues encountered.

### Spec Review

For each spec file:
- Ensure requirements are atomic (one SHALL per requirement)
- Ensure requirements are testable
- Ensure requirements are implementation-agnostic
- Ensure requirements use appropriate EARS patterns

### Critique

When specs are complete, suggest the user run `/sdd/tools/critique specs`:

- Checks for completeness and contradictions
- Identifies missing edge cases
- Validates requirements are well-formed

### Completion

When they explicitly approve the specs:

1. Update state.md: `## Phase Status: complete`, clear `## Notes`
2. Suggest running `/sdd/discovery <name>`

Do not log completion in `## Pending` (that section is for unresolved blockers/decisions only).
