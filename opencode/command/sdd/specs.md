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

!`cat changes/$1/state.md 2>/dev/null || echo "State file not found"`

!`cat changes/$1/proposal.md 2>/dev/null || echo "No proposal found"`

### Entry Check

Apply state entry check logic from `sdd-state-management` skill.

If lane is not `full`, redirect user to appropriate command.

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
