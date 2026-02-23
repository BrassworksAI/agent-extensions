# references/ Guide (When Needed)

Use `references/` only when it reduces token load or confusion.

## When to add references

Add `references/` when one or more are true:

- `SKILL.md` is getting long (approaching ~500 lines) or hard to scan.
- The skill spans multiple domains (e.g., finance vs sales) and you want to avoid loading irrelevant details.
- You have a large body of examples that shouldn't be loaded unless requested.
- You need a detailed API reference, schema list, or company policy document.

Default: keep a skill single-file unless there's a clear benefit.

## How to structure references

- Keep references one level deep from `SKILL.md`.
- Avoid chains like `SKILL.md` -> `advanced.md` -> `details.md`.
- Name files by purpose, not by number:
  - Good: `references/auth-patterns.md`, `references/api-conventions.md`
  - Bad: `references/doc1.md`, `references/misc.md`

## Linking from SKILL.md

In `SKILL.md`, link to references using a relative path so the agent can always locate the file:

```markdown
See `/references/api-conventions.md`.
```

Skills are plain `SKILL.md` files; do not use template variables.

## Long reference files

If a reference file is longer than ~100 lines, add a small table of contents at the top so the agent can quickly navigate.
