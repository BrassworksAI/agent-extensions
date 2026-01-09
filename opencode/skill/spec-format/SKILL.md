---
name: spec-format
description: How to write SDD specifications - EARS syntax, delta format, requirement structure
---

# Spec Format

This skill covers how to write and structure SDD specifications.

## Spec File Structure

Canonical specs (post-finish) may live in a domain taxonomy under `specs/` at the repository root.

Change-set specs for a given change set live under `changes/<name>/specs/` and may be nested by domain/subdomain.

```
specs/
  <domain>/
    <subdomain>/
      <capability>.md
```

## Capability Spec Format

Each capability spec follows this structure:

```markdown
# <Capability Name>

## Overview

Brief description of what this capability does and why it exists.

## Access (optional)

Include this section only when access constraints are meaningful for the capability.

> Authenticated

or

> Public

Projects define their own access levels in a platform-level spec (e.g., `auth-gating.md`).

## Dependencies (optional)

Optional section listing dependencies on other specs or external systems.

## Requirements

Requirements as a bulleted list. Optional grouping headings (### Group Name) for organization.

### <Optional Group>

- The system SHALL <requirement in active voice>.
- WHEN <trigger> the system SHALL <action>.

### <Another Group>

- WHILE <state> the system SHALL <action>.
```

### Topic headings (groups)

Under `## Requirements`, group headings are simply whatever the next heading level down is:

- In `kind: new` specs, groups/topics are typically `### <Topic>`
- In `kind: delta` specs, groups/topics live under `### ADDED/MODIFIED/REMOVED`, so topics are `#### <Topic>`

`Entry` and `Exit` are just common topic names for user-facing or system-facing flows (they are not special section types).

Example (new spec):

```markdown
## Requirements

### Entry

- WHEN a user navigates to <destination> the system SHALL <action>.

### Exit

- WHEN <action> completes the system SHALL redirect to <destination>.
```

### Surface-Agnostic Language

Use action-oriented but surface-agnostic language in Entry/Exit:

| Prefer | Avoid |
|--------|-------|
| "navigates to" | "clicks link to" |
| "initiates sign-out" | "clicks sign-out button" |
| "requests deletion" | "submits the delete form" |

The spec describes *what the user intends to do*, not *how they physically do it*.

## EARS Syntax

Requirements use EARS (Easy Approach to Requirements Syntax) patterns:

| Pattern | Template | Use When |
|---------|----------|----------|
| Ubiquitous | The system SHALL `<action>`. | Fundamental system properties, always true |
| Event-driven | WHEN `<trigger>` the system SHALL `<action>`. | Response to a specific event |
| State-driven | WHILE `<state>` the system SHALL `<action>`. | Behavior during a particular state |
| Unwanted behavior | IF `<condition>` THEN the system SHALL `<action>`. | Handling errors, failures, edge cases |
| Optional feature | WHERE `<feature>` is present the system SHALL `<action>`. | Behavior tied to optional features |
| Complex | WHEN `<trigger>` IF `<condition>` THEN the system SHALL `<action>`. | Combining patterns |

### Examples

- The system SHALL validate all user input before processing.
- WHEN the user clicks submit the system SHALL save the form data.
- WHILE in maintenance mode the system SHALL reject new connections.
- IF the database connection fails THEN the system SHALL retry with exponential backoff.
- WHERE two-factor auth is enabled the system SHALL require a verification code.

## Change Set Specs (`changes/<name>/specs/`)

All change set specs live in `changes/<name>/specs/` and MUST include YAML frontmatter:

```markdown
---
kind: new | delta
---
```

### Optional frontmatter fields

Change-set specs MAY include additional YAML frontmatter fields when helpful. This is intentionally extensible for future tooling.

#### `depends_on` (optional)

Use `depends_on` to express spec-to-spec dependencies in a machine-readable form:

```yaml
---
kind: new
depends_on:
  - specs/<domain>/<...>/<capability>.md
---
```

Prefer listing canonical spec paths under `specs/` even when the dependency is being modified in the current change set.

### `kind: new`

Use `kind: new` when the spec does not exist in canonical yet. Write it like a normal spec.
Do NOT use delta markers or delta sections.

### `kind: delta`

Use `kind: delta` when the spec modifies an existing canonical spec.

`kind: delta` specs use section buckets (NOT blockquote markers):

- `### ADDED`
- `### MODIFIED`
- `### REMOVED`

Under each bucket, use `#### <Topic>` headings.

For MODIFIED topics, requirements MUST be represented as adjacent Before/After blocks:

```markdown
### MODIFIED

#### <Topic>

**Before:**
- The system SHALL <old text>.

**After:**
- The system SHALL <new text>.
```

For ADDED topics:

```markdown
### ADDED

#### <Topic>

- The system SHALL <new requirement>.
```

For REMOVED topics:

```markdown
### REMOVED

#### <Topic>

- The system SHALL <requirement being removed>.

**Reason:** <Why this requirement is being removed>
```

### Where ADDED/MODIFIED/REMOVED may appear

Only use `### ADDED/MODIFIED/REMOVED` inside `## Requirements` and `## Access` (if `## Access` is present).

- `## Access` is optional; omit it for foundational/mechanism specs (e.g., navigation, rate limiting).
- When `## Access` is present in a delta spec, do not add `####` topics under it; put the Before/After or bullets directly under the corresponding `###` bucket.

## RFC 2119 Keywords

Requirements use RFC 2119 language:

- **SHALL** - Absolute requirement
- **SHALL NOT** - Absolute prohibition
- **SHOULD** - Recommended but not required
- **SHOULD NOT** - Discouraged but not prohibited
- **MAY** - Optional

Prefer SHALL for most requirements. Use SHOULD/MAY sparingly.

## Writing Good Requirements

1. **One requirement per bullet** - Don't combine multiple requirements
2. **Active voice** - "The system SHALL validate" not "Validation shall be performed"
3. **Testable** - Requirements should be verifiable
4. **Implementation-agnostic** - Describe WHAT, not HOW
5. **No ambiguity** - Avoid words like "appropriate", "reasonable", "user-friendly"
6. **Use EARS patterns** - Pick the appropriate pattern for the requirement type

## Referencing Requirements

Tasks, plans, and reconciliation reference requirements by quoting the exact EARS line from the spec.

Example in tasks.md:
```markdown
**Requirements:**
- "WHEN the user clicks submit the system SHALL save the form data."
- "IF validation fails THEN the system SHALL display error messages."
```

## Flat vs Grouped Delta Specs

**Flat:** One delta file per modified capability (default)
```
changes/<name>/specs/auth/login.md
changes/<name>/specs/auth/logout.md
```

**Grouped:** Multiple related changes in one file (for tightly coupled changes)
```
changes/<name>/specs/auth/session-management.md  # covers login + logout + refresh
```

Use grouped only when changes are truly interdependent and reviewing separately would lose context.
