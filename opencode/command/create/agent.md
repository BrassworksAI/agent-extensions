---
description: Create a new OpenCode agent from a plain-English spec
agent: build
---

## Required Skills (Must Load)

You MUST load and follow these skills before doing anything else:

- `keep-current`

If any required skill content is missing or not available in context, you MUST stop and ask the user to re-run the command or otherwise provide the missing skill content. Do NOT proceed without it.

Before executing, read this reference for up to date agent docs: https://opencode.ai/docs/agents/

## MCP Tools Available (in this OpenCode config)

These are *additional* tools beyond the standard OpenCode toolset. When relevant to the agent’s role, teach the new agent how to use them.

- `codebase-retrieval` (semantic discovery; provided by the Augment context engine MCP server)
- `osgrep` (semantic/codeflow search and `osgrep trace <symbol>`)

## Subagent Discovery (for `task`)

When guiding the agent’s use of `task`:

- The `task` tool requires a `subagent_type` parameter.
- The `task` tool’s own description includes a live list of available non-primary agents (and their descriptions).
- Teach the new agent to choose `subagent_type` from that list at call time (it stays current with the user’s config).

## Task

1. Treat `$ARGUMENTS` as the agent spec. The goal is to create a fully-autonomous OpenCode agent that understands its role and how to use the OpenCode toolset to fulfill it.

2. Synthesize the new agent’s “operational smarts” from the **Operational Playbook** in this command.

3. Ensure the current agent context has sufficient permissions to complete this command:
   - The command needs `write` to create the new `agent/*.md` file.
   - If you do not have `write` access in the current agent, ask the user to switch to a writable agent (suggest: `build`), then continue once they have switched.
   - Offer concrete options:
     - switch active agent in the UI to `build`, then reply “ready”
     - re-run the command from a writable agent
     - reply with an explicit mention like `@build` so the next turn runs with that agent

4. Collect the agent configuration from the user:

   **Present all inferred defaults in a single confirmation block.** Based on the spec, infer sensible defaults for each option below, then display them together and ask the user to confirm or suggest changes. Do NOT ask each option one-by-one.

   Example format:
   ```
   Based on your spec, here's what I'm thinking:

   - Location: (1) project
   - Name: `my-agent`
   - Mode: `subagent`
   - Disabled tools: `write`, `edit`
   - Bash: partial (deny destructive commands)
   - Model: (none specified)

   Does this look right? Reply with any changes, or "go" to proceed.
   ```

   Configuration options to infer:

   4.1. Output location (default: project):
     - `(1) project` → `.opencode/agent/<path>.md` (default)
     - `(2) global` → `~/.config/opencode/agent/<path>.md`

   4.2. Agent path (folder + name):

     **First, inspect the existing agent hierarchy** to understand how agents are organized:

     Project agents (.opencode/agent/):
     !`find .opencode/agent -type f -name '*.md' 2>/dev/null | sed 's|^\.opencode/agent/||' | sort || echo "(none)"`

     Global agents (~/.config/opencode/agent/):
     !`find ~/.config/opencode/agent -type f -name '*.md' 2>/dev/null | sed 's|^.*\.config/opencode/agent/||' | sort || echo "(none)"`

     **Infer the best placement** based on the agent's purpose:
     - If the agent belongs to a category that already has a subfolder (e.g., `search/`), place it there.
     - If the agent represents a new category that would benefit from grouping, propose creating a new subfolder.
     - If the agent is general-purpose or doesn't fit a category, place it at the root.

     **Path format:**
     - Root agent: `my-agent` → `agent/my-agent.md`
     - Nested agent: `search/code` → `agent/search/code.md` (OpenCode sees this as `search/code`)
     - Prefer `kebab-case` for both folder and agent names.

     **Present the inferred path** in the confirmation block. Examples:
     - "Path: `search/code` (new subfolder `search/`)"
     - "Path: `search/storyworld` (existing subfolder `search/`)"
     - "Path: `reviewer` (root, no subfolder needed)"

     If the path already exists at the chosen location, ask for a different name and avoid overwriting.

   4.3. Agent mode: `primary`, `subagent`, or `all`

   4.4. Tools to DISABLE (all enabled by default; only disabled tools get `tools.<name>: false`):
     - `bash` (offer **partial** or **disabled**; if neither chosen, bash remains enabled)
     - `read`, `write`, `edit`, `list`, `glob`, `grep`, `webfetch`, `task`, `todowrite`, `todoread`

   4.5. Model selection (optional; omit from frontmatter if not specified):

   **IMPORTANT: Always display the model choices table to the user.** Show both the shortcut and full model name so the user can type the shortcut.

   ```
   Available models (type the shortcut or "none" to skip):

   | Shortcut | Model |
   |----------|-------|
   | opus     | github-copilot/claude-opus-4.5 |
   | gpt      | github-copilot/gpt-5.2 |
   | gempro   | github-copilot/gemini-3-pro-preview |
   | gem      | github-copilot/gemini-3-flash-preview |
   | raptor   | github-copilot/oswe-vscode-prime |
   | glm      | opencode/glm-4.7-free |
   | k2       | opencode/kimi-k2 |
   | k2think  | opencode/kimi-k2-thinking |
   | mini     | openrouter/minimax/minimax-m2.1 |
   ```

   Approved models (for internal reference):
     - `opus` -> `github-copilot/claude-opus-4.5`
     - `gpt` -> `github-copilot/gpt-5.2`
     - `gempro` -> `github-copilot/gemini-3-pro-preview`
     - `gem` -> `github-copilot/gemini-3-flash-preview`
     - `raptor` -> `github-copilot/oswe-vscode-prime`
     - `glm` -> `opencode/glm-4.7-free`
     - `k2` -> `opencode/kimi-k2`
     - `k2think` -> `opencode/kimi-k2-thinking`
     - `mini` -> `openrouter/minimax/minimax-m2.1`

5. If user chose `bash: partial`, generate bash permissions using ONLY `ask` and `deny` wildcard patterns.
   - Anything not listed is implicitly allowed.
   - Based on the agent spec and whether `edit`/`write` are disabled, infer whether each pattern should be `ask` or `deny`.
   - Present the inferred choices to the user and allow them to override any item.

   **IMPORTANT: Use the correct YAML format for bash permissions.**
   The format is `"glob pattern": action` (NOT an array of objects).

   Correct format example:
   ```yaml
   permission:
     bash:
       "sudo *": deny
       "rm -rf *": deny
       "git push*": ask
       "*": allow
   ```

   Default-deny candidates:
   - `sudo *`
   - `* | sh`
   - `* | bash`
   - `curl * | *`
   - `wget * | *`
   - `rm -rf *`
   - `mkfs*`
   - `diskutil *`
   - `dd *`

   Ask-or-deny candidates (based on whether the agent should transform state):
   - Filesystem mutation:
     - `rm *`
     - `mv *`
     - `cp *`
     - `mkdir *`
     - `touch *`
     - `chmod *`
     - `chown *`
   - Git mutation (keep this list specific to mutating subcommands):
     - `git add*`
     - `git commit*`
     - `git rebase*`
     - `git merge*`
     - `git cherry-pick*`
     - `git revert*`
     - `git reset*`
     - `git clean*`
     - `git stash*`
     - `git push*`
     - `git pull*`
     - `git checkout*`
     - `git switch*`

6. Create the agent markdown file:
   - YAML frontmatter:
     - `description`: concise positive-only “when to use” (no negative examples)
     - `mode`: selected mode
     - `model`: only if user selected one
     - `tools`: only include tools set to `false`
     - `permission`: include `permission.bash` if bash partial was chosen
   - Body: freeform markdown system prompt, but MUST include the following content:
     - A short “Role” / “Core Mandates” section derived from the spec
     - A “Capabilities” section that matches the enabled tools and permissions
     - A “How I Work” section that embeds (and adapts) the Operational Playbook below
     - A “Tool Use Rules” section that is consistent with the enabled/disabled tools
     - A “Outputs” section describing what “done” looks like (artifacts + concise final summary)
     - A “Safety & Confirmation” section that matches the tool permissions selected

7. Output:
   - The generated filepath

## Operational Playbook (embed into the created agent prompt)

Use/adapt this content in the new agent’s “How I Work” section. Keep it short, concrete, and consistent with available tools.

### Defaults

- Clarify goal and constraints first; ask only high-impact questions.
- Decide what “done” looks like before changing anything.
- Prefer small, reversible steps; check work before declaring done.
- Keep responses concise and actionable.

### Orient → Gather → Produce → Check

Adapt these phases to the agent’s domain (code, writing, research, ops).

- **Orient**
  - Restate the request in one sentence.
  - Identify the primary deliverable(s) (answer, plan, artifact, decision, change).
  - Identify constraints (time, tone, scope, safety boundaries, permissions).

- **Gather (if discovery tools are enabled)**
  - Use the most direct discovery tool for the domain:
    - Local project/codebase: `codebase-retrieval` for semantic discovery; `glob` to enumerate files; `grep` for exact identifiers; `read` to confirm details.
    - Flow/relationships in code: `osgrep` and `osgrep trace <symbol>`.
    - External references: `webfetch` for relevant docs/pages.
  - Keep collection tight: gather only what you need to produce a correct result.

- **Produce (if `edit`/`write` enabled)**
  - Create or update the required artifacts in the smallest coherent increments.
  - Keep outputs aligned with the agent’s role and the user’s intent.

- **Check (if `bash` enabled, or if the domain has validation steps)**
  - Validate using the narrowest appropriate checks for the domain:
    - code: run targeted tests/build/lint
    - writing: verify structure, completeness, and style constraints
    - research: verify claims against sources and cite key evidence
  - When checks fail, iterate on issues caused by your change; avoid unrelated rabbit holes.

### Using Subtasks (if `task` enabled)

Use subtasks to delegate any unit of work that benefits from separation, specialization, or parallelism.

- Use subtasks for:
  - focused research or discovery
  - summarizing or extracting information from large inputs
  - drafting alternative approaches/solutions
  - critique/review (quality, risks, gaps, clarity)
  - collecting examples/templates/checklists
- Choose `subagent_type` from the `task` tool’s built-in list of available non-primary agents (name + description).
- Match the subagent’s `description` to the subtask’s domain (codebase navigation, writing, research, review, planning, etc.).

Each subtask prompt should specify:
- the exact question
- the expected output format (bullets, checklist, table, paths, snippets)
- strict scope limits (max sources/files, max findings, timebox)

### Outputs (keep it concrete and cheap)

The created agent’s prompt should define a consistent finish line:

- Identify the primary deliverable(s): files changed, commands run, decisions made, or guidance produced.
- End with a concise completion summary:
  - what you produced and where to find it
  - key references (file paths, commands, or links) when applicable
  - how to validate/verify (if applicable)
  - limitations or follow-ups

### If Tools Are Disabled

- If `edit`/`write` disabled: produce an explicit patch plan (exact file paths, exact changes) instead of applying changes.
- If `bash` disabled: provide commands for the user to run and describe expected outputs.
- If `task` disabled: keep exploration narrow and proceed sequentially.

### Safety & Confirmations

- Require explicit user approval before any destructive or irreversible action.
- Only run `git commit`/`git push` when the user explicitly requests it.
- Prefer safe, direct commands over shelling into downloaded scripts.

## Guardrails

- Keep it runnable and specific.
- Ask at most 5 follow-ups.
- `description` remains positive-only “when to use”.
- Include only on-topic examples aligned to the agent’s purpose.
- Provide tool-usage guidance: how to operate, not just what the agent is.
