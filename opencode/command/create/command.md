---
description: Create a new OpenCode command from a plain-English description
agent: build
---

## Required Skills (Must Load)

You MUST load and follow these skills before doing anything else:

- `bun-shell-commands`

If any required skill content is missing or not available in context, you MUST stop and ask the user to re-run the command or otherwise provide the missing skill content. Do NOT proceed without them.

Go read the OpenCode commands docs before proceeding: https://opencode.ai/docs/commands/

## Task

1. Treat `$ARGUMENTS` as the command spec
2. Ask user for output path:
   - `(1) project` → `.opencode/command/<name>.md`
   - `(2) global` → `~/.config/opencode/command/<name>.md`
3. Choose a short name based on the spec. Prefer one work but two word kebab-case is okay if necessary.
4. Present the user with a list of agents to choose from for the agent frontmatter. Also give them an `auto` option and you choose based on the command spec
   - Project agents: !`ls -1 .opencode/agent 2>/dev/null | sed 's/\.md$//' || true`
   - Global agents: !`ls -1 ~/.config/opencode/agent 2>/dev/null | sed 's/\.md$//' || true`
5. Determine if the command should include shell scripts (using !`command` syntax). If yes, use the bun-shell-commands skill to learn how to write them. Shell scripts are useful for deterministic context that doesn't require LLM tool calls.
6. Write the new command markdown (docs format): frontmatter (`description`, `agent`), agent spec and goal in markdown
7. Output: filepath, selected agent (brief why), and an example invocation.

## Guardrails

- Keep it specific and runnable.
- Ask at most 5 follow-ups if required.
- Use shell scripts (!`command`) for deterministic operations that don't need LLM reasoning (e.g., listing files, checking git status, getting environment info). Reference @skill/bun-shell-commands/SKILL.md for proper syntax and patterns.

## When to Use Inline Shell Scripts

Use !`command` to inject deterministic context when:

- **Reading predictable files**: Inject contents of all files in a folder that the agent will need to analyze
  ```markdown
  !`for f in $(ls src/*.ts); do echo "=== $f ===" && cat "$f"; done`
  ```
- **Git operations**: Get branch name, status, or diff upfront
  ```markdown
  !`git branch --show-current && git status --short`
  ```
- **Project metadata**: Read package.json, Cargo.toml, or similar config files
  ```markdown
  !`cat package.json`
  ```
- **File discovery**: List all files matching a pattern before processing
  ```markdown
  !`find . -name "*.test.ts" | head -20`
  ```
- **Environment info**: Get OS, shell, or other runtime details
  ```markdown
  !`uname -a && echo $SHELL`
  ```
- **Dependency inspection**: List installed packages or tools
  ```markdown
  !`npm list --depth=0 2>/dev/null || echo "No package.json"`
  ```

**Key principle**: If the agent will inevitably need this information and there's no decision-making required to get it, inject it upfront with a shell script.
