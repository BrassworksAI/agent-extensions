---
name: skill-creator
description: Collaboratively design and author agent skills with correct frontmatter, naming, and placement conventions.
---

# Skill Creator

## Optional: Quick fit check

Good fit if the user wants to:

- Create a new skill
- Improve an existing `SKILL.md`
- Decide where a skill should live (project/global/custom)

Not needed if the user is asking for a normal code change, debugging help, or a one-off question.

If this skill isn’t needed, ignore it and continue.

## Outcomes

When used, this skill helps produce:

- A new skill folder containing a valid `SKILL.md`
- A clear, discovery-friendly `description`
- Optional `references/` and `scripts/` only when needed

## How skills work (progressive disclosure)

This skill follows the Agent Skills model:

- **Discovery**: agents load only `name` + `description` for all skills.
- **Activation**: when this skill is selected, the agent reads `SKILL.md`.
- **Execution**: the agent reads `references/` files or runs `scripts/` only when needed.

Design new skills the same way: keep the main file high-signal, and push deep details into files that are read only on the branch that needs them.

## Core rules

- Skills are discovered by **frontmatter only**: `name` + `description`.
- Skill `name` must match the directory name and pass naming constraints.
- Default to a single-file skill.
- Add `references/` only when it reduces token load or avoids confusion.
- Add `scripts/` only for deterministic, mechanical, or fragile steps.

### Minimal spec constraints (inline)

- `name`:
  - 1–64 characters
  - lowercase letters/numbers/hyphens only
  - no leading/trailing `-`, no `--`
  - must match the skill directory name
- Recall: `description` is used for discovery. It should say what the skill does and when to use it, and include concrete trigger keywords.

## Workflow

### 1) Gather requirements (collaborative)

This is a collaborative workflow: ask questions first, then draft, then confirm.

Ask for:

- Skill name (kebab-case)
- A single-sentence `description` that makes the discovery decision easy
- What the skill should reliably produce (outputs/artifacts)
- What questions the agent must ask up front (inputs)
- Guardrails (what to avoid, what not to assume)

Do not generate files until the user confirms the name, description, and placement.

### 2) Choose placement (three options)

Ask the user to choose where the new skill should live:

1. **Project local (Recommended)**: `.<tool>/skills/<skill-name>/SKILL.md`
2. **Global**: `~/.<tool>/skills/<skill-name>/SKILL.md`
3. **Other (absolute path)**: user provides an absolute path

Rules for **Other**:

- If the provided path ends with `SKILL.md`, treat it as the exact file path.
- Otherwise treat it as a directory path and place the skill at `<path>/<name>/SKILL.md`.
- If the path isn’t absolute, ask again for an absolute path.

### 3) Validate before writing

Before creating files, verify:

- `name` satisfies the constraints above
- The directory name will exactly match `name`
- `description` is 1–1024 characters and describes when to use the skill

### 4) Author `SKILL.md`

Write the new `SKILL.md` directly. Use this structure as a starting point (adapt as needed):

```markdown
---
name: <skill-name>
description: <One sentence describing when this skill should be used>
---

# <Skill Title>

## Collaboration

- Ask the user for missing inputs.
- Summarize intended outputs and file changes.
- Wait for confirmation before creating or editing files.

## Inputs I need

- <question 1>
- <question 2>

## Workflow

1. <Step 1>
2. <Step 2>
3. <Step 3>

## Guardrails

- <constraints>
```

Guidelines:

- Prefer concise, step-based workflows.
- Use gates/decision points when paths branch.
- Keep the skill focused on what the user’s environment and repo actually require.

### 4a) References and scripts (relative paths)

Skills are plain `SKILL.md` files (no template rendering). When referencing bundled files, always use relative paths from the skill root:

- References: `/references/<doc>.md`
- Scripts: `/scripts/<script>.mjs` (or `.ts` if needed)

### 5) Add `references/` (branch only)

Default: do not add `references/`.

Add `references/` only when it keeps `SKILL.md` smaller and more usable (e.g. approaching ~500 lines, domain-specific docs, long examples).

If you decide you need references, first read:

- `/references/references-guide.md`

When linking to bundled reference files from `SKILL.md`, use a relative path:

```markdown
See `/references/<doc>.md`.
```

### 6) Add `scripts/` (branch only, Node-only)

Default: do not add `scripts/`.

Add scripts only for deterministic, mechanical steps that benefit from automation (validation, scaffolding, transforms).

If you decide you need scripts:

1. Read script conventions:
   - `/references/scripts-overview.md`

2. Use `.mjs` scripts and the Node standard library only (`node:*` imports).

3. Validate scripts parse before you rely on them:
   - Node: `node -c scripts/<script>.mjs`

## How to write a good `description`

A strong `description` makes the discovery decision obvious.

Patterns that work well:

- Start with an action verb: “Design…”, “Generate…”, “Review…”, “Refactor…”, “Draft…”, “Validate…”
- Include the artifact produced: “...a `SKILL.md`…”, “...a migration plan…”, “...a checklist…”
- Include discovery triggers: what user phrases or contexts should cause this skill to activate

Examples:

- “Creates a scoped `SKILL.md` for consistent API endpoint changes in this repo. Use when adding or modifying API routes or handlers.”
- “Drafts agent skills for repeatable workflows, keeping `SKILL.md` concise and using `references/` only for deep details.”
