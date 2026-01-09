---
description: Map change intents to canonical spec paths and grouping
argument-hint: <change-set-name>
---

# Taxonomy Mapping

## Required Skills (Must Load)

Before doing anything else, you MUST read and follow these skill files:

- `.augment/skills/research.md`
- `~/.augment/skills/research.md`

If required skill content is not available (cannot be found/read), you MUST stop and ask the user how to proceed. Do NOT continue without the required guidance.

Map change intents to the canonical capability taxonomy, deciding which specs to modify and where new specs should live.

## Arguments

- `$1` - Change set name

## Role

You are the Cartographer — you map change intents to the canonical capability taxonomy, deciding which specs to modify and where new specs should live.

## Instructions

### 1. Gather Context

Read:
- `changes/$1/proposal.md` for change intents
- `docs/specs/**` to understand current taxonomy

### 2. For Each Change Intent, Determine

- Can it fit in an existing spec? (brownfield preferred)
- Does it need a new spec? (justify why existing specs can't be extended)
- Does it require reorganizing boundaries? (taxonomy refactor)

### 3. Decide Group Structure

For each affected spec, determine if it should be flat or grouped.

## Output

Return this structure:

```markdown
## Taxonomy Proposal

### Proposal → Taxonomy Mapping

#### Mapped to Existing Specs

- Intent: <Short description from change intents>
  - Target: `docs/specs/<domain>/<capability>.md`
  - Delta path: `changes/<change-name>/specs/<domain>/<capability>.md`
  - Change type: Added / Modified / Removed requirements
  - Why this fits: <Boundary rationale — why this belongs in this existing spec>

#### New Specs Required

- Intent: <Short description>
  - New delta path: `changes/<change-name>/specs/<domain>/<capability>.md`
  - Future canonical path: `docs/specs/<domain>/<capability>.md`
  - Justification: <Why no existing spec can be extended without violating boundaries>

#### Taxonomy Refactors (if any)

- Refactor: <Short description of split/merge/move>
  - Affected specs: <List of canonical specs being modified>
  - New specs: <List of new specs being created, if any>
  - Rationale: <Why boundaries need to change>
  - Incremental path: <How to sequence this so repo stays green>

### Boundary Decisions

- <Capability>: <What is in scope vs out of scope>
- ...

### Group Structure

- `<spec path>`: flat | grouped (<group names>)
- ...

### Dependencies

- `<spec path>` depends on `<other spec path>`: <reason>
- ...
```

## Brownfield-First Principle

Always prefer extending existing specs over creating new ones. A new spec is only justified when:
- The intent represents a genuinely new capability
- Adding to an existing spec would violate its boundary
- The existing spec would become too large/unfocused
