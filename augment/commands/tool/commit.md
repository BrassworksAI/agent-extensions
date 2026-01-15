---
description: Smart commits with semantic grouping and conventional commit format
---

# Smart Commit

Create conventional commits from staged/unstaged changes. `$ARGUMENTS` provides optional context.

## Current Git State

### Branch
!`git branch --show-current`

### Status
!`git status --short`

### Staged Changes
!`git diff --cached --stat`

### Unstaged Changes
!`git diff --stat`

### Full Diff
!`git diff HEAD`

### Recent Commits (style reference)
!`git log --oneline -10`

## Rules

**Format:** `<type>(<scope>): <description>`

**Types:** `feat` | `fix` | `refactor` | `docs` | `test` | `chore`

**Scope:** Module affected (`auth`, `api`, `ui`, `db`, etc.) - match existing patterns

**Message style:**
- Single sentence only, no body
- Describe what the system does now (not what you did)
- Used in changelogs - write for developers

**Examples:**
- `feat(auth): support OAuth2 refresh tokens`
- `fix(api): prevent duplicate webhook deliveries`
- `refactor(db): consolidate connection pooling logic`

**Avoid:** `feat(auth): added code for tokens` · `fix: fixed bug` · `update files`

## Execute

1. Group changes by concern - one commit per feature/fix
2. For each commit:
   - All files: `git add -A && git commit -m "<type>(<scope>): <description>"` (or `git commit -a -m "..."`)
   - Specific files: `git add <files> && git commit -m "<type>(<scope>): <description>"`
3. Verify: `git log --oneline -5`
