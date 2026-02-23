---
name: merge-change-specs
description: Use this is you are trying to merge change-set specs from `changes/<name>/specs/` into canonical `specs/`
---

# Merge Change Specs

Merges change-set specs into canonical specs with deterministic output.

## Usage

```bash
node ./scripts/merge-change-specs.mjs --change <name> [--dry-run]
```

**Flags:**

- `--change <name>` (required): The change set name (directory under `changes/`)
- `--dry-run` (optional): Preview changes without writing files

## Workflow

1. **Dry run first** - Always preview changes before applying:

   ```bash
   node ./scripts/merge-change-specs.mjs --change auth-refresh --dry-run
   ```

2. **Review the JSON output** - Check created/modified/skipped counts and file list

3. **Execute if ready** - Run without `--dry-run` to apply:

   ```bash
   node ./scripts/merge-change-specs.mjs --change auth-refresh
   ```

## Output

The script outputs a JSON object:

```json
{
  "change": "auth-refresh",
  "dryRun": true,
  "counts": {
    "created": 1,
    "modified": 2,
    "skipped": 0
  },
  "created": ["specs/auth/login.md"],
  "modified": ["specs/auth/session.md", "specs/auth/logout.md"],
  "skipped": []
}
```

If validation fails or paths are unsafe, the script exits non-zero with errors to stderr.

## Spec Format

Change-set specs must follow the format defined in the `spec-format` skill:

- YAML frontmatter with `kind: new` or `kind: delta`
- Valid markdown structure per spec-format validation

The script automatically validates each spec using the `spec-format` validator.
