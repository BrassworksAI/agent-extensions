---
description: Create a new Augment agent from a plain-English spec
argument-hint: <agent-description>
---

# Create Agent

## Required Skills (Must Load)

Before doing anything else, you MUST read and follow these skill files:

- `.augment/skills/keep-current.md`
- `~/.augment/skills/keep-current.md`

If required skill content is not available (cannot be found/read), you MUST stop and ask the user how to proceed. Do NOT continue without the required guidance.

Create a new Augment CLI agent from a plain-English specification.

## Arguments

- `$ARGUMENTS` - Plain-English description of what the agent should do

## Task

1. Treat `$ARGUMENTS` as the agent spec. The goal is to create a fully-autonomous Augment agent that understands its role and how to use the available toolset to fulfill it.

2. Collect the agent configuration from the user:

   **Present all inferred defaults in a single confirmation block.** Based on the spec, infer sensible defaults for each option below, then display them together and ask the user to confirm or suggest changes. Do NOT ask each option one-by-one.

   Example format:
   ```
   Based on your spec, here's what I'm thinking:

   - Location: (1) project
   - Name: `my-agent`

   Does this look right? Reply with any changes, or "go" to proceed.
   ```

   Configuration options to infer:

   2.1. Output location (default: project):
     - `(1) project` → `.augment/agents/<name>.md`
     - `(2) global` → `~/.augment/agents/<name>.md`

   2.2. Agent name:
     - Infer a kebab-case name from the spec
     - Check if the name already exists at the chosen location
     - If it exists, suggest an alternative name

3. Create the agent markdown file:

   **Frontmatter** (only `description` is supported):
   ```yaml
   ---
   description: <concise "when to use" description>
   ---
   ```

   **Body**: Freeform markdown system prompt including:
   - A short "Role" / "Core Mandates" section derived from the spec
   - A "How I Work" section with operational guidance
   - An "Outputs" section describing what "done" looks like

4. Output:
   - The generated filepath
   - Example invocation: `@agent-name`

## Agent File Location

Augment agents go in:
- **Project-local**: `.augment/agents/<name>.md`
- **Global**: `~/.augment/agents/<name>.md`

## Operational Playbook (embed into the created agent prompt)

Use/adapt this content in the new agent's "How I Work" section:

### Defaults

- Clarify goal and constraints first; ask only high-impact questions.
- Decide what "done" looks like before changing anything.
- Prefer small, reversible steps; check work before declaring done.
- Keep responses concise and actionable.

### Orient → Gather → Produce → Check

- **Orient**
  - Restate the request in one sentence.
  - Identify the primary deliverable(s).
  - Identify constraints (scope, safety boundaries).

- **Gather**
  - Use discovery tools to understand context:
    - `codebase-retrieval` for semantic search
    - `grep-search` for exact identifiers
    - `view` to read file contents
  - Keep collection tight: gather only what you need.

- **Produce**
  - Create or update artifacts in small coherent increments.
  - Keep outputs aligned with the agent's role.

- **Check**
  - Validate using appropriate checks for the domain.
  - When checks fail, iterate on issues caused by your change.

### Outputs

- Identify the primary deliverable(s): files changed, commands run, guidance produced.
- End with a concise completion summary:
  - what you produced and where to find it
  - key references (file paths, commands)
  - how to validate/verify (if applicable)
  - limitations or follow-ups

## Guardrails

- Keep it runnable and specific.
- Ask at most 5 follow-ups.
- `description` remains positive-only "when to use".
