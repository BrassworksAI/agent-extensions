---
name: spec-format
description: How to write SDD specifications - EARS syntax, delta format, requirement structure
---

# Spec Format

This skill covers how to write and structure SDD specifications.

## Spec File Structure

Specs live in `specs/` at the repository root, organized by domain:

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

## Access

Optional section declaring access level for capabilities that have access control.
Use a blockquote marker to indicate the access level:

> Authenticated

or

> Public

Projects define their own access levels in a platform-level spec (e.g., `auth-gating.md`).
This keeps individual specs clean while still being explicit about access requirements.

## Dependencies

Optional section listing dependencies on other specs or external systems.

## Entry

Optional section describing how users/systems arrive at this capability.
Use for capabilities with clear entry points (navigation, triggers, conditions).

> **ADDED**

- WHEN a user navigates to <destination> the system SHALL <action>.

## Requirements

Requirements as a bulleted list. Optional grouping headings (### Group Name) for organization.

### <Optional Group>

- The system SHALL <requirement in active voice>.
- WHEN <trigger> the system SHALL <action>.

### <Another Group>

- WHILE <state> the system SHALL <action>.

## Exit

Optional section describing how the capability completes or where users end up.
Use for capabilities with clear exit points (redirects, returns, completions).

> **ADDED**

- WHEN <action> completes the system SHALL redirect to <destination>.
- IF <condition> THEN the system SHALL <alternative exit>.
```

### When to Use Entry/Exit

| Spec Type | Entry | Exit |
|-----------|-------|------|
| User-facing flows (sign-in, deletion) | ✓ | ✓ |
| Compositions (navigation, dashboard) | — | — |
| Programmatic capabilities (read, provision) | ✓ (trigger) | ✓ (returns) |
| Static content pages (terms, privacy) | ✓ | — |

Omit sections that don't apply. Not every spec needs Entry/Exit.

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

## Delta Spec Format

When modifying specs during a change set, create delta specs in `changes/<name>/specs/`:

```
changes/<name>/specs/
  <domain>/
    <capability>.md
```

Delta specs use markers to show changes:

### Adding New Requirements

```markdown
### <Group>

> **ADDED**

- The system SHALL <new requirement>.
- WHEN <trigger> the system SHALL <action>.
```

### Modifying an Existing Requirement

```markdown
### <Group>

> **MODIFIED**

**Before:**
- The system SHALL <old text>.

**After:**
- The system SHALL <new text>.
```

### Removing a Requirement

```markdown
### <Group>

> **REMOVED**

- The system SHALL <requirement being removed>.

**Reason:** <Why this requirement is being removed>
```

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
