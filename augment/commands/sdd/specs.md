---
description: Write delta specifications for change
argument-hint: <change-set-name>
---

# Specs

Write delta specifications for the change set.

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

Specs live in `specs/` at the repository root, organized by domain:

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

#### Delta Spec Markers

When creating delta specs in `changes/<name>/specs/`, use markers:

**Adding:**
```markdown
> **ADDED**

- The system SHALL <new requirement>.
```

**Modifying:**
```markdown
> **MODIFIED**

**Before:**
- The system SHALL <old text>.

**After:**
- The system SHALL <new text>.
```

**Removing:**
```markdown
> **REMOVED**

- The system SHALL <requirement being removed>.

**Reason:** <Why this requirement is being removed>
```

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
