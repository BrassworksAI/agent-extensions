# scripts/ Overview (When Needed)

Scripts are for deterministic, mechanical steps that benefit from automation.

Default: do not bundle scripts.

## When scripts are worth it

Good fit:

- Validation: naming rules, required sections, frontmatter parsing.
- Scaffolding: creating a common folder/file layout.
- Transforms: mechanical edits that are easy to get subtly wrong.
- Repeatability: the same action is performed frequently.

Bad fit:

- The task is mostly judgment, writing, or design decisions.
- The workflow is one-off and low-risk.
- The script would need lots of dependencies or network access.

Rule of thumb:

- Put judgment and policy in `SKILL.md`.
- Put mechanical steps in `scripts/`.

## Script conventions

- Put scripts under `scripts/`.
- Prefer `.mjs` and Node standard library only.
- Keep scripts single-purpose and small.
- Validate CLI args early. Print usage and exit non-zero on invalid usage.

## Safety

- Require explicit paths as arguments (avoid implicit global writes).
- Print what will be created/modified before writing.
- Avoid deleting files by default.
- If deletion is necessary, require an explicit `--force` flag.

## Output patterns

Default output should be concise and actionable:

- What the script did
- What it changed (paths)
- What to do next

Example:

```text
Created:
- .opencode/skill/my-skill/SKILL.md

Next:
- Edit SKILL.md description and workflow
```

Optional: support `--json` if another tool will consume output.

Rules for `--json`:

- JSON must be the only output.
- Keep fields stable.

## Validation loop

Use a tight feedback loop:

1. Run the script
2. Review output and warnings
3. Fix issues
4. Re-run until clean
