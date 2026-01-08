---
name: sdd/init
description: Initialize a new SDD change set
agent: sdd/plan
---

<skill>sdd-state-management</skill>

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

   ## Problem

   (What problem are we solving?)

   ## Goals

   (What does success look like?)

   ## Non-Goals

   (What are we explicitly NOT doing?)

   ## Approach

   (High-level approach - freeform)
   ```
6. **Report**: Confirm creation and suggest next step (`/sdd/brainstorm` or `/sdd/proposal`)
