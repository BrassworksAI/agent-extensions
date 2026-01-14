# Delta Merge Rules (merge-change-specs)

This merger assumes change-set specs follow `spec-format`.

## Paths

- Change-set specs live at `changes/<name>/specs/**/*.md`
- Canonical specs live at `specs/**/*.md`
- Canonical path is computed by stripping the `changes/<name>/specs/` prefix.

## Frontmatter

Change-set specs MUST begin with YAML frontmatter:

```yaml
---
kind: new | delta
---
```

Frontmatter is never copied into canonical.

## `kind: new`

- Canonical file content becomes the change-set spec body (frontmatter removed).

## `kind: delta`

Delta files are treated as merge instructions against the canonical spec.

- Canonical spec MUST exist or the merge fails.

### Requirements buckets

Buckets under the delta file’s `## Requirements` section are interpreted as operations:

- `### REMOVED`: delete the exact bullet blocks from canonical, and delete the topic section if it becomes empty.
- `### MODIFIED`: replace exact Before blocks with After blocks.
- `### ADDED`: insert bullet blocks into the matching canonical topic (`### <Topic>`) within `## Requirements`.

### MODIFIED topics

Each topic MUST encode an adjacent Before/After block:

```markdown
### MODIFIED

#### <Topic>

**Before:**
- ...exact text as it currently exists...

**After:**
- ...replacement text...
```

Merge behavior:

- The tool searches the canonical spec for each exact **Before** block.
- It replaces it with the exact **After** block.
- Labels (`**Before:**`, `**After:**`) are not written into canonical.
