---
description: Write change-set specifications for change
argument-hint: <change-set-name>
---

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

Before writing specs, research the codebase:

1. **Delegate to @librarian** to understand:
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

### Writing Delta Specs

Create specs in `changes/<name>/specs/` following the spec format guidance below.

#### Spec File Structure

Canonical specs (post-finish) may live in a domain taxonomy under `specs/` at the repository root.

Change-set specs for this change are authored under `changes/<name>/specs/` and may be nested by domain/subdomain.

```
specs/
  <domain>/
    <subdomain>/
      <capability>.md
```

#### Capability Spec Format

```markdown
# <Capability Name>

## Overview

Brief description of what this capability does and why it exists.

## Requirements

Requirements as a bulleted list. Optional grouping headings for organization.

### <Optional Group>

- The system SHALL <requirement in active voice>.
- WHEN <trigger> the system SHALL <action>.
```

#### EARS Syntax

Requirements use EARS (Easy Approach to Requirements Syntax) patterns:

| Pattern | Template | Use When |
|---------|----------|----------|
| Ubiquitous | The system SHALL `<action>`. | Fundamental system properties |
| Event-driven | WHEN `<trigger>` the system SHALL `<action>`. | Response to a specific event |
| State-driven | WHILE `<state>` the system SHALL `<action>`. | Behavior during a particular state |
| Unwanted behavior | IF `<condition>` THEN the system SHALL `<action>`. | Handling errors, edge cases |
| Complex | WHEN `<trigger>` IF `<condition>` THEN the system SHALL `<action>`. | Combining patterns |

#### Change Set Spec Kinds

All specs under `changes/<name>/specs/` MUST include YAML frontmatter:

```markdown
---
kind: new | delta
---
```

- Use `kind: new` for brand new capabilities; write the spec like a normal spec (no delta markers).
- Use `kind: delta` for edits to existing canonical specs.

#### Delta Spec Format

For `kind: delta`, use section buckets (NOT blockquote markers):

```markdown
## Requirements

### ADDED

#### <Topic>
- The system SHALL <new requirement>.

### MODIFIED

#### <Topic>

**Before:**
- The system SHALL <old text>.

**After:**
- The system SHALL <new text>.

### REMOVED

#### <Topic>
- The system SHALL <removed requirement>.

**Reason:** <Why>
```

Only use `### ADDED/MODIFIED/REMOVED` inside `## Requirements` and `## Access` (if present).
### Identifying Capabilities

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
3. Suggest running `/sdd:discovery <name>`

Don't advance until the user clearly signals approval. Questions, feedback, or acknowledgments don't count as approval.
