---
description: Advance to the next SDD phase with artifact checks
---

# SDD Next

Advance the active change set to the next phase only after verifying current-phase artifacts exist.

## Required Skills

- `spec-driven-development` (state management, phase sequencing, guardrails)

## Inputs

> [!IMPORTANT]
> Resolve the active change set from workspace state first. Ask the user only when multiple change sets are present.

## Instructions

1. Load `spec-driven-development` and run `ae sdd status [name]`.

2. Identify lane and current phase.

3. Verify required artifacts for the current phase before advancing:

   - **full/proposal**: `changes/<name>/proposal.md` exists.
   - **full/specs**: `changes/<name>/specs/` exists and contains at least one `.md` spec.
   - **full/discovery**: `changes/<name>/thoughts/` exists and contains at least one discovery note.
   - **full/tasks**: `changes/<name>/tasks.toml` exists and contains at least one task entry.
   - **full/plan**: `changes/<name>/plans/` exists and contains the current task plan.
   - **full/implement**: `changes/<name>/tasks.toml` exists and task progress is 100%.
   - **full/reconcile**: `changes/<name>/reconciliation.md` exists.
   - **vibe/context** and **bug/triage**: `changes/<name>/context.md` exists.
   - **vibe/plan** and **bug/plan**: `changes/<name>/plan.md` exists.
   - **vibe/reconcile** and **bug/reconcile**: `changes/<name>/reconciliation.md` exists.

4. If artifacts are missing, do not advance. Report exactly what is missing and suggest the command to produce/fix it.

5. If artifacts are present, advance phase by running:

   ```bash
   ae sdd phase complete --next [name]
   ```

6. Report the transition result and the new phase from status output.

## Success Criteria

- Phase progression happens only after artifact checks pass.
- Missing artifacts block progression with actionable guidance.
- Successful progression is explicit and visible in status output.
