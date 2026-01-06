---
name: sdd/specs
description: Write change-set specifications for change
agent: sdd/plan
---

<skill>sdd-state-management</skill>
<skill>spec-format</skill>
<skill>research</skill>

# Specs

Write change-set specifications for the change set (`kind: new` and `kind: delta`).

## Arguments

- `$ARGUMENTS` - Change set name

## Instructions

### Setup

1. Read `changes/<name>/state.md` - verify phase is `specs` and lane is `full`
2. Read `changes/<name>/proposal.md` for context
3. List existing `specs/` structure to understand current taxonomy

### Research Phase

Before writing specs, use the `research` skill:

1. **Research to understand**:
   - Current spec structure and taxonomy
   - Related existing capabilities
   - How similar things are specified

2. **Build context** for spec writing:
   - What specs already exist in related areas
   - What naming conventions are used

### Taxonomy Mapping

With research in hand, suggest the user run `/sdd/tools/taxonomy-map <name>`:

- Determines where new capabilities should live in the spec hierarchy
- Recommends brownfield (existing specs) vs greenfield (new specs)
- Provides boundary decisions and group structure

### Writing Change Set Specs

Create specs in `changes/<name>/specs/` following the `spec-format` skill. Specs may be nested by domain/subdomain under that folder (e.g. `changes/<name>/specs/auth/login.md`).
Remember that change set specs have YAML frontmatter `kind: new | delta`.

1. **Identify capabilities** needed from the proposal
2. **Determine paths** using taxonomy mapping guidance
3. **Write requirements** using EARS syntax

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

Work through specs collaboratively with the user. When they explicitly approve:

1. Log approval in state.md under `## Pending`:
   ```
   None - Specs approved: [list of spec files written]
   ```
2. Update state.md phase to `discovery`
3. Suggest running `/sdd/discovery <name>`

Don't advance until the user clearly signals approval. Questions, feedback, or acknowledgments don't count as approval.
