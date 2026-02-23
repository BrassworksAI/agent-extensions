# Spec management notes (CLI + repo workflow)

## Goals

- Deterministic, diffable spec metadata (overview, requirements, access level, etc.)
- Human- and LLM-friendly authoring
- Ability to compute a spec diff for a change set and/or branch
- Centralized specs repo (avoid per-product repo sprawl)

## Core approach (recommended)

**Keep specs in Markdown with YAML frontmatter** and enforce a strict schema + normalization rules.

- **Frontmatter**: flat, stable keys only (no nested maps); arrays sorted and normalized on write.
- **Body**: freeform markdown for rationale, examples, diagrams.
- **Determinism**: CLI validates required fields, canonicalizes ordering, and rewrites frontmatter in a stable order.

Example (frontmatter only):

```yaml
---
id: spec-auth-roles
title: Role-based access control
overview: Add RBAC to admin-only endpoints
access_level: internal
requirements:
  - rbac model + roles matrix
  - endpoint enforcement
  - audit logging
owners:
  - team-auth
repo:
  url: https://github.com/acme/product
  base: main
change_set: rbac-rollout
---
```

## Spec repository layout (centralized)

```text
specs/
  change-sets/
    rbac-rollout/
      state.toml
      tasks.toml
      specs/
        spec-auth-roles.md
        spec-api-gates.md
```

- Keep **state/tasks** in TOML (already used by `ae sdd`).
- Keep **spec content** in Markdown with YAML frontmatter.

## Diff strategy

### 1) Local git diff (preferred)

Use git to compute deterministic deltas without API drift:

- `git diff --name-status <base>...<head> -- specs/`
- `git diff <base>...<head> -- specs/change-sets/<name>`
- If desired, strip or normalize frontmatter before diffing to hide re-ordering noise.

This is stable and avoids remote API inconsistencies. The CLI can `git fetch` and diff locally.

### 2) Remote diff via provider API (optional)

Fallback for users who do not have the repo checked out locally.

- GitHub: PR compare endpoint or `gh pr diff`.
- GitLab/Bitbucket: compare API endpoints.

This should be *secondary* due to auth and API rate limits.

## Branch and change-set management

If specs live in a central repo, the CLI can manage worktrees:

- `ae sdd repo init` → config for remote + base branch
- `ae sdd checkout <change-set>` → create branch, track in state
- `ae sdd sync` → fetch, rebase/merge base, update state

Git worktrees can allow multiple change sets in parallel without branch switching.

## Proposed CLI surface (minimal)

```text
ae sdd spec init <change-set> <spec-id>
ae sdd spec list [change-set]
ae sdd spec validate [change-set]
ae sdd spec normalize [change-set]
ae sdd spec diff [change-set] [--base main]
ae sdd repo init --url <repo> --base main
ae sdd repo checkout <change-set>
ae sdd repo sync [change-set]
```

## Determinism rules (important for clean diffs)

- Canonical frontmatter key order (fixed list)
- Arrays sorted and de-duplicated
- Line endings normalized (LF)
- No nested maps unless explicitly required
- Required fields enforced (e.g., `id`, `title`, `overview`, `access_level`)

## Why this fits LLM + human workflows

- YAML frontmatter is easy to edit and machine-validate.
- Markdown body preserves narrative context and examples.
- Deterministic frontmatter output removes diff noise.
- Git diffs remain the single source of truth for “what changed.”

## Open questions

- Should `repo.url` + `repo.base` live in each spec or be inherited at change-set level?
- Do we want a spec index file per change set for fast listing and search?
- Should CLI support an “offline mode” that skips remote fetch entirely?
