---
description: Fast-track bug investigation and fix initialization
---

# SDD Bug

## Required Skills

- `spec-driven-development` (state management, lane/phase flow)
- `research`

## Inputs

> [!IMPORTANT]
> Resolve required inputs from the workspace first; only ask the user when resolution is ambiguous or missing (for example, multiple change sets exist).

- **Bug Context**: Ask for error messages, reproduction steps, or a description of what is failing.
- **Change Set Name**: Resolve by checking existing change sets. If only one exists, proceed with it; if multiple exist, ask whether this is new or related to an existing one. If none exist, ask the user to confirm or provide a name for the new change set (e.g., `fix-login-error`).

## Instructions

1. **Resolution**: Check for existing change sets. If multiple exist, ask the user if this is a new issue or related to an existing one.

2. **Triage**: Determine if this is an **Actual Bug** (code fails to meet existing specs/intent) or a **Behavioral Change** (new capability requested). If it's a change, redirect the user to the full lane (`/sdd/init`).

3. **Initialize** (if new): Create the change set using the bug lane:

   ```bash
   ae sdd init <name> --lane bug
   ```

   This creates `state.toml` with lane=bug, phase=plan, status=in_progress.

4. **Research**: Use the `research` skill to locate the problem, trace the execution path, and identify the root cause. Compare findings against existing specs.

5. **Document**: Create `changes/<name>/context.md` with:
   - Problem description
   - Expected behavior
   - Root cause (from research)
   - High-level fix approach

6. Update state notes with findings: `ae sdd notes set "<summary>"`

7. **Next Steps**: Instruct the user to run `/sdd/plan` to start the fix.

## Success Criteria

- Bug triaged correctly (actual bug vs behavioral change)
- Change set initialized with correct `bug` lane configuration
- Root cause documented in `context.md`
- Phase ready for planning

## Usage Examples

### Do: Triage before fixing

"This appears to be a regression where the login button is disabled on mobile. I'll initialize a bug lane change set `fix-mobile-login` to research the root cause."

### Don't: Assume CLI arguments

Avoid starting work before asking the user for the bug details or the desired change set name.
