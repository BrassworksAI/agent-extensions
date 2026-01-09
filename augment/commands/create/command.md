---
description: Create a new Augment command from a plain-English description
argument-hint: <command-description>
---

# Create Command

## Required Skills (Must Load)

Before doing anything else, you MUST read and follow these skill files:

- `.augment/skills/bun-shell-commands.md`
- `~/.augment/skills/bun-shell-commands.md`

If required skill content is not available (cannot be found/read), you MUST stop and ask the user how to proceed. Do NOT continue without the required guidance.

Create a new Augment CLI command from a plain-English description.

## Arguments

- `$ARGUMENTS` - Plain-English description of what the command should do

## Task

1. Treat `$ARGUMENTS` as the command spec

2. Ask user for output path:
   - `(1) project` → `.augment/commands/<name>.md`
   - `(2) global` → `~/.augment/commands/<name>.md`

3. Choose a short name based on the spec:
   - Prefer one word but two-word kebab-case is okay if necessary
   - Can use nested paths for organization: `sdd/my-command` creates `commands/sdd/my-command.md`

4. Collect configuration in a single confirmation block:

   ```
   Based on your spec, here's what I'm thinking:

   - Location: (1) project
   - Name: `my-command`
   - Argument hint: <brief-arg-description>

   Does this look right? Reply with any changes, or "go" to proceed.
   ```

5. Write the new command markdown:

   **Frontmatter**:
   ```yaml
   ---
   description: <concise description of what the command does>
   argument-hint: <optional hint for arguments>
   ---
   ```

   **Body**: Markdown with:
   - Clear instructions for what the command does
   - Steps to execute
   - Any prompts or questions to ask the user
   - Expected outputs

6. Output:
   - The generated filepath
   - An example invocation: `/my-command`

## Command File Location

Augment commands go in:
- **Project-local**: `.augment/commands/<name>.md`
- **Global**: `~/.augment/commands/<name>.md`

## Command Frontmatter Options

Augment commands support these frontmatter fields:

| Field | Description |
|-------|-------------|
| `description` | Brief description shown in command list |
| `argument-hint` | Hint shown after command name for arguments |
| `model` | Override the model for this command |

## Writing Good Commands

### Structure

```markdown
---
description: What this command does
argument-hint: <args>
---

# Command Name

Brief intro.

## Arguments

- `$ARGUMENTS` - Description of expected arguments
- `$1` - First positional argument (if using positional)

## Instructions

1. Step one
2. Step two
3. Step three

## Output

What the command produces.
```

### Tips

- Keep commands focused on a single task
- Use clear, imperative instructions
- Define what success looks like
- Use `$ARGUMENTS` for the full argument string, `$1`, `$2` etc. for positional args

## Guardrails

- Keep it specific and runnable.
- Ask at most 5 follow-ups if required.
- Commands should be self-contained with clear instructions.
