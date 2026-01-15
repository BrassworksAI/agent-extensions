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

1. **Project local (Recommended)**: `{{{SKILL_LOCAL_PATH}}}/<skill-name>/SKILL.md`
2. **Global**: `{{{SKILL_GLOBAL_PATH}}}/<skill-name>/SKILL.md`
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

### 4a) Template files (`.tmpl.md`)

If your skill references scripts or supporting files using path variables, name the file with a `.tmpl.md` suffix instead of `.md`.

**When to use `.tmpl.md`:**

- The skill bundles scripts that need to be invoked with a known install path
- The skill bundles `references/` files and needs stable cross-install paths
- Any file that contains `{{{SKILL_INSTALL_PATH}}}`, `{{{SKILL_LOCAL_PATH}}}`, or `{{{SKILL_GLOBAL_PATH}}}`

**How it works:**

- Files named `*.tmpl.md` are rendered at install time with variables substituted
- The output is placed in a cache directory (e.g., `.cache/skills-rendered/opencode/`)
- Symlinks point to the rendered output, not the source template
- Non-template files (`*.md`) are symlinked directly to the source

**Available template variables:**

| Variable | Description | Example (OpenCode global) |
|----------|-------------|---------------------------|
| `{{{SKILL_INSTALL_PATH}}}` | Path to installed skill location | `$HOME/.config/opencode/skill` |
| `{{{SKILL_LOCAL_PATH}}}` | Project-local skill directory | `.opencode/skill` |
| `{{{SKILL_GLOBAL_PATH}}}` | Global skill directory | `~/.config/opencode/skill` |

**Conflict detection:**

If both `SKILL.md` and `SKILL.tmpl.md` exist in a skill directory, the install will fail with an error. Choose one or the other.

### 5) Add `references/` (branch only)

Default: do not add `references/`.

Add `references/` only when it keeps `SKILL.md` smaller and more usable (e.g. approaching ~500 lines, domain-specific docs, long examples).

If you decide you need references, first read:

- `{{{SKILL_INSTALL_PATH}}}/skill-creator/references/references-guide.md`

When linking to bundled reference files from `SKILL.md`, use an install-stable path:

```markdown
See `{{{SKILL_INSTALL_PATH}}}/my-skill/references/<doc>.md`.
```

### 6) Add `scripts/` (branch only, Node-first)

Default: do not add `scripts/`.

Add scripts only for deterministic, mechanical steps that benefit from automation (validation, scaffolding, transforms).

If you decide you need scripts:

1. Read script conventions:
   - `{{{SKILL_INSTALL_PATH}}}/skill-creator/references/scripts-overview.md`

2. Pick runtime (Node-first):
   - `{{{SKILL_INSTALL_PATH}}}/skill-creator/references/scripts-runtime-node-first.md`

3. Use `.mjs` scripts and the Node standard library only (`node:*` imports).

4. Validate scripts parse/compile before you rely on them:

   - Node: `node -c scripts/<script>.mjs`
   - Bun (fallback runner): `bun --check scripts/<script>.mjs`

## How to write a good `description`

A strong `description` makes the discovery decision obvious.

Patterns that work well:

- Start with an action verb: “Design…”, “Generate…”, “Review…”, “Refactor…”, “Draft…”, “Validate…”
- Include the artifact produced: “...a `SKILL.md`…”, “...a migration plan…”, “...a checklist…”
- Include discovery triggers: what user phrases or contexts should cause this skill to activate

Examples:

- “Creates a scoped `SKILL.md` for consistent API endpoint changes in this repo. Use when adding or modifying API routes or handlers.”
- “Drafts agent skills for repeatable workflows, keeping `SKILL.md` concise and using `references/` only for deep details.”
