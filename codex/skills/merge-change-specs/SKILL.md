---
name: merge-change-specs
description: Merge change-set specs in `changes/<name>/specs/` into canonical `specs/` deterministically.
---

# Merge Change Specs

## Inputs

- Change set name (the `<name>` in `changes/<name>/`).
- Optional: `--dry-run` to preview changes.

## Workflow

1. Validate the change set exists and contains `specs/**/*.md`.
2. For each change-set spec:
   - Validate the markdown format using the `spec-format` validator script.
   - Parse YAML frontmatter and determine `kind: new|delta`.
   - Compute the canonical spec path by stripping `changes/<name>/specs/` and prefixing with `specs/`.
3. Apply changes deterministically:
   - `kind: new`: write canonical file as the change-set body (frontmatter removed).
   - `kind: delta`: patch the canonical spec by applying the delta file’s `### ADDED`, `### MODIFIED`, `### REMOVED` buckets.
4. Emit a deterministic JSON summary (created/modified/skipped) and fail fast on errors.

## Guardrails

- Fail fast if a `kind: delta` spec targets a missing canonical file.
- Do not reorder unrelated canonical content; only apply targeted edits.
- Keep ordering deterministic: stable traversal + stable output formatting.

## Commands

- Run merge: `node codex/skills/merge-change-specs/scripts/merge-change-specs.mjs --change <name>`
- Dry run: `node codex/skills/merge-change-specs/scripts/merge-change-specs.mjs --change <name> --dry-run`

Note: This script imports the validator from `codex/skills/spec-format/scripts/validate-change-spec.mjs`.

## References

- Merge semantics: `references/delta-merge-rules.md`
- Spec format + validator: `codex/skills/spec-format/SKILL.md`