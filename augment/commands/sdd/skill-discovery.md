---
description: Debug skill discovery and loading paths
argument-hint: [skill-name...]
---

# Skill Discovery (Debug)

This command helps troubleshoot skill loading by explicitly attempting to locate and read skills from expected paths.

## Arguments

- `$ARGUMENTS` - Optional list of skill filenames or basenames to probe (space-separated). If omitted, probe a small default set.

## Instructions

### 1) Determine which skills to probe

- If `$ARGUMENTS` is provided, treat each token as a skill name.
- If no arguments are provided, use this default list:
  - `sdd-state-management`
  - `research`
  - `spec-format`
  - `architecture-fit-check`
  - `architecture-workshop`

### 2) For each skill, try to load in this order

For each skill name `X`, attempt to read the first file that exists:

1. `.augment/skills/X.md`
2. `~/.augment/skills/X.md`

If a file is missing, note it and continue to the next location.

### 3) Report results

For each skill, report:

- Whether a project-local file was found
- Whether a global file was found
- Which file was actually loaded (if any)
- Any read errors encountered

### 4) Provide a short diagnosis

Based on what you found, explain:

- Which paths are being missed
- Whether the issue is missing files vs. discovery order
- Any obvious fixes (e.g., wrong folder, wrong filename, missing `.md`)

## Output

Provide a concise, human-readable summary. Use a short table if helpful.

## Example

```
/sdd:skill-discovery research spec-format
```
