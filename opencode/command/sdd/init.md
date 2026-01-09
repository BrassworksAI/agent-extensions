---
name: sdd/init
description: Initialize a new SDD change set
agent: sdd/plan
---

# Required Skills (Must Load)

You MUST load and follow these skills before doing anything else:

- `sdd-state-management`

If any required skill content is missing or not available in context, you MUST stop and ask the user to re-run the command or otherwise provide the missing skill content. Do NOT proceed without them.

# Initialize Change Set

Create a new SDD change set with the given name.

## Arguments

- `$ARGUMENTS` - Name for the change set (kebab-case)

## Instructions

1. **Validate name**: Ensure name is kebab-case, no spaces, lowercase
2. **Check for conflicts**: Verify `changes/<name>/` doesn't already exist
3. **Create structure**:
   ```
   changes/<name>/
     state.md
     proposal.md
   ```
4. **Initialize state.md**:
   ```markdown
   # SDD State: <name>

   ## Lane

   (not yet selected)

   ## Phase

   ideation

   ## Phase Status

   in_progress

   ## Pending

    - Select lane (full/vibe/bug) during proposal

   ## Notes
   ```
5. **Initialize proposal.md** with empty template:
   ```markdown
    # Proposal: <name>

    ## Context

    (Drop any initial notes, requirements, or information here to help inform the proposal)

    ## Problem

    (What problem are we solving?)

    ## Goals

    (What does success look like?)

    ## Approach

    (High-level approach - freeform)
    ```
6. **Report**: Confirm creation and suggest next step (`/sdd/brainstorm` or `/sdd/proposal`)
