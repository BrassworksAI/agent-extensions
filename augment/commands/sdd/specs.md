---
description: Write change-set specifications for change
argument-hint: <change-set-name>
---

# Specs

Write change-set specifications for change set (`kind: new` and `kind: delta`).

## Arguments

- `$ARGUMENTS` - Change set name

## Instructions

> **SDD Process**: Read `.augment/skills/sdd-state-management.md` for state management guidance.

> **Research**: When needed, delegate to `@librarian` for codebase context. See `.augment/skills/research.md` (project-local) or `~/.augment/skills/research.md` (global) for guidance.

> **Spec Format**: Use guidance from `.augment/skills/spec-format.md` (project-local) or `~/.augment/skills/spec-format.md` (global) for EARS syntax and structure.

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
