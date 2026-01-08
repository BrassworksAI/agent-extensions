---
description: Smart commits with semantic grouping and conventional commit format
agent: sdd/build
---

# Smart Commit

Analyze staged and unstaged changes to create well-structured conventional commits.

## Current Git State

### Branch
!`git branch --show-current`

### Status
!`git status --short`

### Staged Changes
!`git diff --cached --stat`

### Unstaged Changes
!`git diff --stat`

### Full Diff (for semantic analysis)
!`git diff HEAD`

### Recent Commits (for style reference)
!`git log --oneline -10`

## Your Task

Analyze the changes above and create commits following these rules:

### Conventional Commit Format

```
<type>(<scope>): <description>
```

### Allowed Types
- `feat` - New feature or capability
- `fix` - Bug fix
- `refactor` - Code restructuring without behavior change
- `docs` - Documentation only
- `test` - Adding or updating tests
- `chore` - Build, tooling, or maintenance tasks

### Scope Guidelines
- Scope should reflect the **system/module** affected (e.g., `auth`, `api`, `ui`, `db`)
- Keep scopes consistent with existing commit history
- Omit scope only for cross-cutting changes

### Commit Message Style
- **Single sentence only** - no additional text, body, or clarifying details
- Write in terms of what the **system now does** (not what you did)
- Focus on the **why** and **impact**, not implementation details
- Messages will be used in changelogs, so write for end-users/developers

**Good examples:**
- `feat(auth): support OAuth2 refresh tokens`
- `fix(api): prevent duplicate webhook deliveries`
- `refactor(db): consolidate connection pooling logic`

**Bad examples:**
- `feat(auth): added code for tokens` (implementation-focused)
- `fix: fixed bug` (no context)
- `update files` (meaningless)

## Execution Steps

1. **Analyze the diff semantically** - Group related changes by feature, fix, or module
2. **Determine commit strategy:**
    - If changes are cohesive (single purpose) → one commit
    - If changes span multiple concerns → suggest a series of commits with specific files for each
3. **For each commit:**
    - **Determine optimal git command:**
      - Check if the commit will include ALL files (staged + unstaged)
      - If yes, use shorthand (`git add -A && git commit` or `git commit -a`)
      - If no, stage specific files (`git add <files>`)
    - **Execute:**
      - All files: `git add -A && git commit -m "<type>(<scope>): <description>"` or `git commit -a -m "..."`
      - Specific files: `git add <files> && git commit -m "<type>(<scope>): <description>"`
4. **Verify** - Run `git log --oneline -5` to confirm commits were created

If the user provided arguments, treat them as additional context: $ARGUMENTS
