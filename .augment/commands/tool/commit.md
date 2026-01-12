---
description: Smart commits with semantic grouping and conventional commit format
argument-hint: [additional context]
---

# Smart Commit

Analyze staged and unstaged changes to create well-structured conventional commits.

## Arguments

- `$ARGUMENTS` - Optional additional context (e.g., "focus on API changes first", "skip test files")

## Instructions

### 1. Gather Git State

Run these commands to understand the current state:

```bash
# Current branch
git branch --show-current

# Status overview
git status --short

# Staged changes
git diff --cached --stat

# Unstaged changes
git diff --stat

# Full diff for semantic analysis
git diff HEAD

# Recent commits for style reference
git log --oneline -10
```

### 2. Analyze Changes Semantically

Review the diff output and identify:

- **Cohesive changes**: Single purpose, single commit
- **Multi-concern changes**: Multiple features/fixes that should be separate commits

Group related changes by:
- Feature or capability
- Bug fix
- Module or system affected

### 3. Conventional Commit Format

Use this format for all commits:

```
<type>(<scope>): <description>
```

#### Allowed Types

| Type | When to use |
|------|-------------|
| `feat` | New feature or capability |
| `fix` | Bug fix |
| `refactor` | Code restructuring without behavior change |
| `docs` | Documentation only |
| `test` | Adding or updating tests |
| `chore` | Build, tooling, or maintenance tasks |

#### Scope Guidelines

- Scope should reflect the **system/module** affected (e.g., `auth`, `api`, `ui`, `db`)
- Keep scopes consistent with existing commit history (check `git log`)
- Omit scope only for cross-cutting changes

### 4. Commit Message Style

- **Single sentence only** - no additional text, body, or clarifying details
- Write in terms of what the **system now does**, not what you did
- Messages will be used in changelogs

**Good examples:**
- `feat(auth): support OAuth2 refresh tokens`
- `fix(api): prevent duplicate webhook deliveries`
- `refactor(db): consolidate connection pooling logic`

**Bad examples:**
- `feat(auth): added code for tokens` (implementation-focused)
- `fix: fixed bug` (no context)
- `update files` (meaningless)

### 5. Execute Commits

**Determine optimal git command:**
- Check if the commit will include ALL files (staged + unstaged)
- If yes, use the shorthand approach (see below)
- If no, stage specific files manually

**Single commit** (cohesive changes):

**Option 1: All files being added (optimized)**
```bash
# If all untracked files + all modified files should be in the commit
git add -A && git commit -m "<type>(<scope>): <description>"

# OR if only tracked files changed (no new files)
git commit -a -m "<type>(<scope>): <description>"
```

**Option 2: Specific files only**
```bash
git add <relevant-files>
git commit -m "<type>(<scope>): <description>"
```

**Multiple commits** (multi-concern changes):
1. Stage files for first logical group
2. Commit with appropriate message
3. Repeat for each group

### 6. Verify

Run `git log --oneline -5` to confirm commits were created correctly.

## Output

- Summary of changes analyzed
- Commit(s) created with their messages
- Verification output showing the new commit(s)
