---
description: Write change-set specifications for change
argument-hint: <change-set-name>
---

# Specs

Write change-set specifications for change set (`kind: new` and `kind: delta`).

## Arguments

- `$ARGUMENTS` - Change set name

## Instructions

### Required Skills (Must Load)

Before doing anything else, you MUST read and follow these skill files:

- `.augment/skills/sdd-state-management.md`
- `.augment/skills/spec-format.md`
- `.augment/skills/research.md`

If any of these are missing project-locally, fall back to the user/global locations:

- `~/.augment/skills/sdd-state-management.md`
- `~/.augment/skills/spec-format.md`
- `~/.augment/skills/research.md`

If required skill content is not available (cannot be found/read), you MUST stop and ask the user how to proceed. Do NOT continue without the required guidance.

### Collaborative Discovery (Required)

Spec writing is a **collaborative** process. You MUST NOT immediately start writing specification files.

Before any spec text is written, you MUST:

1. **Demonstrate understanding of the domain**
   - Summarize the user’s goal, actors, workflows, and constraints in your own words.
   - Identify assumptions and explicitly mark them as assumptions.

2. **Demonstrate understanding of the capability taxonomy**
   - Explain how you believe the change maps into the existing capability hierarchy.
   - Call out any unclear boundaries (what is in-scope vs out-of-scope).

3. **Run a user back-and-forth loop**
   - Ask clarifying questions when anything is ambiguous.
   - Propose concrete options when decisions are needed (tradeoffs, boundary choices).
   - WAIT for explicit user answers/approval before proceeding.

Only after the user confirms (or corrects) the understanding above may you proceed to Research/Taxonomy/Spec writing.

### Setup

1. Read `changes/<name>/state.md`
2. Read `changes/<name>/proposal.md`
3. List existing `specs/` structure to understand current taxonomy

### Entry Check

Apply state entry check logic from `.augment/skills/sdd-state-management.md`.

If lane is not `full`, redirect user to appropriate command.

### Research Phase

Before writing specs, delegate to `@librarian`:

1. **Research to understand**:
   - Current spec structure and taxonomy
   - Related existing capabilities
   - How similar things are specified

2. **Build context** for spec writing:
   - What specs already exist in related areas
   - What naming conventions are used

### Taxonomy Mapping

With research in hand, suggest user run `/sdd:tools:taxonomy-map <name>`:

- Determines where new capabilities should live in spec hierarchy
- Recommends brownfield (existing specs) vs greenfield (new specs)
- Provides boundary decisions and group structure

### Writing Change Set Specs

Only after the user explicitly approves the domain summary and taxonomy mapping:

Create specs in `changes/<name>/specs/` following spec format guidance. Specs may be nested by domain/subdomain under that folder (e.g. `changes/<name>/specs/auth/login.md`).
Remember that change set specs have YAML frontmatter `kind: new | delta`.

1. **Identify capabilities** needed from proposal
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

When specs are complete, suggest user run `/sdd:tools:critique specs`:

- Checks for completeness and contradictions
- Identifies missing edge cases
- Validates requirements are well-formed

### Completion

When they explicitly approve specs:

1. Update state.md: `## Phase Status: complete`, clear `## Notes`
2. Suggest running `/sdd:discovery <name>`

Do not log completion in `## Pending` (that section is for unresolved blockers/decisions only).
